package bot

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) HandleClose(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	if t.Account == types.ACCOUNT_B {
		// close or cancel main order
		go func() {
			// get order
			res, err := b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
				// check if orderId is empty string
				if t.OrderID == "" {
					return nil, fmt.Errorf("the orderID is empty string")
				}
				// convert orderID
				orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.OrderID)
				if err != nil {
					b.l.Printf("error converting orderID to binance OrderID: %v", err)
					return nil, err
				}
				return b.bc.GetOrder(ctx, orderID, t.Symbol)
			})
			if err != nil {
				b.l.Printf("error getting order: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			order, ok := res.(*futures.Order)
			if !ok {
				b.l.Printf("unexpected error happened in casting res to *futures.Order: %T", order)
				return
			}

			switch {
			case order.Status == futures.OrderStatusTypeNew:
				// cancel order if not filled
				res, err := b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
					return b.bc.CancelOrder(ctx, order.OrderID, order.Symbol)
				})
				if err != nil {
					b.l.Printf("error cancelling order: %v", err)
					b.handleError(err, t.UserID, t.ID)
					return
				}
				cancelOrder, ok := res.(*futures.CancelOrderResponse)
				if !ok {
					b.l.Printf("unexpected error happened in casting res to *futures.CancelOrderResponse: %T", cancelOrder)
					return
				}
				// message the user
				msg := fmt.Sprintf("The order successfully cancelled\n\nOrderID: %d\nTradeID: %d", cancelOrder.OrderID, t.ID)
				b.MsgChan <- types.BotMessage{
					ChatID: userID,
					MsgStr: msg,
				}
				// update trade
				b.c.UpdateTradeCanceled(t.ID)
			case order.Status == futures.OrderStatusTypeFilled || order.Status == futures.OrderStatusTypePartiallyFilled:
				// close order with market
				res, err := b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
					var side futures.SideType
					if order.Side == futures.SideTypeBuy {
						side = futures.SideTypeSell
					} else {
						side = futures.SideTypeBuy
					}
					return b.bc.CloseOrder(ctx, order.ExecutedQuantity, side, order.Symbol)
				})
				if err != nil {
					b.l.Printf("error closing order: %v", err)
					b.handleError(err, t.UserID, t.ID)
					return
				}
				closeOrder, ok := res.(*futures.CreateOrderResponse)
				if !ok {
					b.l.Printf("unexpected error happened in casting res to *futures.CreateOrderResponse: %T", closeOrder)
					return
				}
				// message the user
				msg := fmt.Sprintf("The order successfully closed\n\nOrderID: %d\nTradeID: %d", closeOrder.OrderID, t.ID)
				b.MsgChan <- types.BotMessage{
					ChatID: userID,
					MsgStr: msg,
				}
				// update trade
				b.c.UpdateTradeClosed(t.ID)
			default:
				b.l.Printf("the state of order is proper for closing: %v", order.Status)
			}
		}()

		// cancel TP order
		go func() {
			// cancel
			res, err := b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
				// check if orderId is empty string
				if t.TPOrderID == "" {
					return nil, fmt.Errorf("the TP orderID is empty string")
				}
				// convert orderID
				orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.TPOrderID)
				if err != nil {
					b.l.Printf("error converting TP orderID to binance OrderID: %v", err)
					return nil, err
				}
				return b.bc.CancelOrder(ctx, orderID, t.Symbol)
			})
			if err != nil {
				b.l.Printf("error cancelling order: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			order, ok := res.(*futures.CancelOrderResponse)
			if !ok {
				b.l.Printf("unexpected error happened in casting res to *futures.CancelOrderResponse: %T", order)
				return
			}
			// message the user
			msg := fmt.Sprintf("The take-profit order successfully cancelled\n\nOrderID: %d\nTradeID: %d", order.OrderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
		}()

		// cancel SL order
		go func() {
			// cancel
			res, err := b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
				// check if orderId is empty string
				if t.SLOrderID == "" {
					return nil, fmt.Errorf("the SL orderID is empty string")
				}
				// convert orderID
				orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.SLOrderID)
				if err != nil {
					b.l.Printf("error converting SL orderID to binance OrderID: %v", err)
					return nil, err
				}
				return b.bc.CancelOrder(ctx, orderID, t.Symbol)
			})
			if err != nil {
				b.l.Printf("error cancelling order: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			order, ok := res.(*futures.CancelOrderResponse)
			if !ok {
				b.l.Printf("unexpected error happened in casting res to *futures.CancelOrderResponse: %T", order)
				return
			}
			// message the user
			msg := fmt.Sprintf("The stop-loss order successfully cancelled\n\nOrderID: %d\nTradeID: %d", order.OrderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
		}()

	} else {
		b.HandleUnderDevelopment(u, ctx)
	}
	return nil
}
