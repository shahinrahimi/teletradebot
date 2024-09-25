package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange/bitmex"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) handleCloseBitmexOrder(ctx context.Context, t *models.Trade, userID int64, action string, orderIDStr string, isFilledClose bool, changeState bool) {

	if orderIDStr == "" {
		return
	}
	// get order
	res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return b.mc.GetOrder(ctx, orderIDStr, t.Symbol)
	})
	if err != nil {
		b.l.Printf("error getting order: %v", err)
		b.handleError(err, userID, t.ID)
		return
	}
	order, ok := res.(*swagger.Order)
	if !ok {
		b.l.Printf("unexpected error happened in casting res to *swagger.Order: %T", order)
		return
	}
	switch order.OrdStatus {
	case string(bitmex.OrderStatusTypeNew):
		_, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.mc.CancelOrder(ctx, orderIDStr)
		})
		if err != nil {
			b.l.Printf("error cancelling order: %v", err)
			b.handleError(err, userID, t.ID)
			return
		}
		// update trade state
		if changeState {
			b.c.UpdateTradeCanceled(t.ID)
		}
		// message the user
		msg := fmt.Sprintf("Order cancelled successfully.\n%s\n\nOrder ID: %s\nTrade ID: %d", strings.ToUpper(action), orderIDStr, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
	case string(bitmex.OrderStatusTypePartiallyFilled), string(bitmex.OrderStatusTypeFilled):
		if isFilledClose {
			_, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return b.mc.CancelOrder(ctx, orderIDStr)
			})
			if err != nil {
				b.l.Printf("error cancelling order: %v", err)
				b.handleError(err, userID, t.ID)
				return
			}
			// update trade state
			if changeState {
				b.c.UpdateTradeCanceled(t.ID)
			}
			// message the user
			msg := fmt.Sprintf("Order cancelled successfully.\n%s\n\nOrder ID: %s\nTrade ID: %d", strings.ToUpper(action), orderIDStr, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
		}
		return
	default:
		b.l.Printf("unknown order status received: %v", order.OrdStatus)
	}
}
