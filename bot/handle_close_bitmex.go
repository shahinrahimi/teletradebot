package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange/bitmex"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleCloseBitmex(u *tgbotapi.Update, ctx context.Context) {

	t, ok := ctx.Value(models.KeyTrade{}).(models.Trade)
	if !ok {
		b.l.Panic("error getting trade from context")
	}
	userID := u.Message.From.ID

	// close or cancel main order
	go func() {
		// get order
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
			// check if orderId is empty string
			if t.OrderID == "" {
				return nil, fmt.Errorf("the orderID is empty string")
			}
			// convert orderID
			return b.mc.GetOrder(ctx, t.OrderID, t.Symbol)
		})
		if err != nil {
			b.l.Printf("error getting order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*swagger.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", order)
			return
		}

		switch order.OrdStatus {
		case bitmex.OrderStatusTypeNew:
			// cancel order
			res, err = b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
				return b.mc.CancelOrder(ctx, order.OrderID)
			})
			if err != nil {
				b.l.Printf("error cancelling order: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			cancelOrder, ok := res.(*swagger.Order)
			if !ok {
				b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", cancelOrder)
				return
			}
			// message user
			msg := fmt.Sprintf("Order has been canceled.\n\nOrderID: %s\nTrade ID: %d", cancelOrder.OrderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
			// update trade
			b.c.UpdateTradeCanceled(t.ID)

		case bitmex.OrderStatusTypePartiallyFilled, bitmex.OrderStatusTypeFilled:
			// close order with market
			res, err = b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
				return b.mc.CloseOrder(ctx, t.Symbol, order.OrderID)
			})
			if err != nil {
				b.l.Printf("error closing order: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			closeOrder, ok := res.(*swagger.Order)
			if !ok {
				b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", closeOrder)
				return
			}
			// message user
			msg := fmt.Sprintf("Order has been closed.\n\nOrderID: %s\nTrade ID: %d", closeOrder.OrderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
			// update trade
			b.c.UpdateTradeClosed(t.ID)
		default:
			b.l.Printf("unexpected order status: %s", order.OrdStatus)

		}
	}()

	// cancel tp order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
			// check if orderId is empty string
			if t.TPOrderID == "" {
				return nil, fmt.Errorf("the TP orderID is empty string")
			}
			return b.mc.CancelOrder(ctx, t.TPOrderID)
		})
		if err != nil {
			b.l.Printf("error cancelling TP order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*swagger.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", order)
			return
		}
		// message the user
		msg := fmt.Sprintf("Take-profit order has been canceled.\n\nOrderID: %s\nTrade ID: %d", order.OrderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}

	}()

	// cancel sl order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, error) {
			// check if orderId is empty string
			if t.SLOrderID == "" {
				return nil, fmt.Errorf("the SL orderID is empty string")
			}
			return b.mc.CancelOrder(ctx, t.SLOrderID)
		})
		if err != nil {
			b.l.Printf("error cancelling SL order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*swagger.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", order)
			return
		}
		// message the user
		msg := fmt.Sprintf("Stop-loss order has been canceled.\n\nOrderID: %s\nTrade ID: %d", order.OrderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()

}
