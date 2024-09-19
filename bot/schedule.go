package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange/binance"
	"github.com/shahinrahimi/teletradebot/exchange/bitmex"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) scheduleOrderReplacement(ctx context.Context, delay time.Duration, orderId int64, t *models.Trade) {
	time.AfterFunc(delay, func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.bc.GetOrder(ctx, orderId, t.Symbol)
		})
		if err != nil {
			b.l.Printf("error getting order by trade: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := (res).(*futures.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to futures.Order: %T", order)
			return
		}

		// replacement if the order is new (not filled not cancelled etc)
		if order.Status == futures.OrderStatusTypeNew {
			b.l.Printf("Order not executed, attempting replacement")

			// cancel order
			res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return b.bc.CancelOrder(ctx, orderId, t.Symbol)
			})
			if err != nil {
				b.l.Printf("error cancelling order: %v", err)
				// TODO add handleError for cases that canceled and the orderId not found
				if apiErr, ok := err.(*common.APIError); ok {
					if apiErr.Code == -2011 {
						// assume the order is already cancelled successfully so it can not found
						b.handleError(err, t.UserID, t.ID)
					} else {
						b.handleError(err, t.UserID, t.ID)
						b.c.UpdateTradeCanceled(t.ID)
						return
					}
				} else {
					b.handleError(err, t.UserID, t.ID)
					b.c.UpdateTradeCanceled(t.ID)
					return
				}
			}

			cancelOrder, ok := (res).(*futures.CancelOrderResponse)
			if !ok {
				b.l.Printf("unexpected error happened in casting error to futures.CancelOrderResponse: %T", cancelOrder)
				//return
			}

			// place new order
			res, po, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, interface{}, error) {
				return b.bc.PlaceTrade(ctx, t)
			})
			if err != nil {
				b.l.Printf("error placing trade: %v", err)
				b.handleError(err, t.UserID, t.ID)
				// change trade state to canceled
				b.c.UpdateTradeCanceled(t.ID)
				return
			}
			createOrder, ok := (res).(*futures.CreateOrderResponse)
			if !ok {
				b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", createOrder)
				return
			}

			preparedOrder, ok := (po).(*binance.PreparedOrder)
			if !ok {
				b.l.Printf("unexpected error happened in casting interface to binance.PreparedOrder: %T", createOrder)
				return
			}

			// schedule
			go b.scheduleOrderReplacement(ctx, preparedOrder.Expiration, createOrder.OrderID, t)

			// update trade order
			orderIdStr := utils.ConvertBinanceOrderID(createOrder.OrderID)
			b.c.UpdateTradePlaced(t.ID, orderIdStr)

			// message the user
			msg := fmt.Sprintf("Order replaced successfully\n\nNewOrder ID: %d\nTrade ID: %d", createOrder.OrderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: t.UserID,
				MsgStr: msg,
			}

		}
	})
}

func (b *Bot) scheduleOrderReplacementBitmex(ctx context.Context, delay time.Duration, orderID string, t *models.Trade) {
	time.AfterFunc(delay, func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.mc.GetOrder(ctx, t.Symbol, orderID)
		})
		if err != nil {
			b.l.Printf("error getting order by trade: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := (res).(*swagger.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", order)
			return
		}
		b.l.Printf("order status: %s", order.OrdStatus)
		if order.OrdStatus == string(bitmex.OrderStatusTypeNew) {
			b.l.Printf("Order not executed, attempting replacement")
			res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return b.mc.CancelOrder(ctx, orderID)
			})
			if err != nil {
				b.l.Printf("error cancelling order: %v", err)
				b.handleError(err, t.UserID, t.ID)
			}

			cancelOrder, ok := (res).(*swagger.Order)
			if !ok {
				b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", cancelOrder)
			}
			// place a new order
			res, po, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, interface{}, error) {
				return b.mc.PlaceTrade(ctx, t)
			})
			if err != nil {
				b.l.Printf("error placing trade: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			preparedOrder, ok := (po).(*bitmex.PreparedOrder)
			if !ok {
				b.l.Printf("unexpected error happened in casting interface to bitmex.PreparedOrder: %T", preparedOrder)
				return
			}

			createOrder, ok := (res).(*swagger.Order)
			if !ok {
				b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", createOrder)
				return
			}

			// schedule
			go b.scheduleOrderReplacementBitmex(ctx, preparedOrder.Expiration, createOrder.OrderID, t)

			// update trade order
			b.c.UpdateTradePlaced(t.ID, createOrder.OrderID)

			// message the user
			msg := fmt.Sprintf("Order replaced successfully\n\nNewOrder ID: %s\nTrade ID: %d", createOrder.OrderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: t.UserID,
				MsgStr: msg,
			}
		}
	})
}
