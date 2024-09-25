package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) scheduleOrderReplacementBinance(ctx context.Context, i *models.Interpreter, t *models.Trade, ex exchange.Exchange) {
	b.c.SetInterpreter(i, t.ID)
	delay := i.CalculateExpiration()
	b.l.Printf("schedule order replacement: delay: %s, TradeID: %d", delay, t.ID)
	time.AfterFunc(delay, func() {
		oe := i.GetOrderExecutionBinance(types.GetOrderExecution, t.OrderID)
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return ex.GetOrder(ctx, oe)
		})
		if err != nil {
			b.l.Printf("error getting order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := (res).(*futures.Order)
		if !ok {
			b.l.Panicf("unexpected error happened in casting error to futures.Order: %T", order)
		}

		switch order.Status {
		case futures.OrderStatusTypeNew, futures.OrderStatusTypeExpired:
			b.l.Printf("Order not executed, attempting replacement, TradeID: %d", t.ID)
			oe := i.GetOrderExecutionBinance(types.CancelOrderExecution, t.OrderID)
			_, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return ex.CancelOrder(ctx, oe)
			})
			if err != nil {
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
			time.Sleep(config.WaitForReplacement)
			// fetch new interpreter
			i, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return ex.FetchInterpreter(ctx, t)
			})
			if err != nil {
				b.l.Printf("error fetching interpreter: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			interpreter, ok := i.(*models.Interpreter)
			if !ok {
				b.l.Panicf("unexpected error happened in casting error to *models.Interpreter: %T", interpreter)
			}
			b.c.SetInterpreter(interpreter, t.ID)
			// place new order
			oe = interpreter.GetOrderExecutionBinance(types.StopPriceExecution, t.OrderID)
			res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return ex.PlaceStopOrder(ctx, oe)
			})
			if err != nil {
				b.l.Printf("error placing stop-order: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			order, ok := (res).(*futures.Order)
			if !ok {
				b.l.Panicf("unexpected error happened in casting error to futures.Order: %T", order)
			}
			orderID := utils.ConvertBinanceOrderID(order.OrderID)
			// update trade
			b.c.UpdateTradeMainOrder(t.ID, orderID)
			// message user
			msg := fmt.Sprintf("Order replaced successfully.\n\nNewOrder ID: %s\nTrade ID: %d", orderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: t.UserID,
				MsgStr: msg,
			}
			// new schedule
			b.scheduleOrderReplacementBinance(ctx, interpreter, t, ex)
		default:
			b.l.Printf("Schedule order replacement canceled due to status: %s, TradeID: %d", order.Status, t.ID)
		}
	})
}
