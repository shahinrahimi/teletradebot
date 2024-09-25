package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) handleCloseBinanceOrder(ctx context.Context, t *models.Trade, userID int64, action string, orderIDStr string, isFilledClose bool, changeState bool) {

	if orderIDStr == "" {
		return
	}
	// convert orderID
	orderID, err := utils.ConvertOrderIDtoBinanceOrderID(orderIDStr)
	if err != nil {
		b.l.Printf("unexpected error converting orderID to binance OrderID: %v", err)
		return
	}
	// get order
	res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return b.bc.GetOrder(ctx, orderID, t.Symbol)
	})
	if err != nil {
		b.l.Printf("error getting order: %v", err)
		b.handleError(err, userID, t.ID)
		return
	}
	order, ok := res.(*futures.Order)
	if !ok {
		b.l.Printf("unexpected error happened in casting res to *futures.Order: %T", order)
		return
	}
	switch order.Status {
	case futures.OrderStatusTypeNew:
		_, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.bc.CancelOrder(ctx, orderID, order.Symbol)
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

		// message user
		msg := fmt.Sprintf("Order closed successfully.\n%s\n\nOrder ID: %s\nTrade ID: %d", strings.ToUpper(action), orderIDStr, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
		return
	case futures.OrderStatusTypeFilled, futures.OrderStatusTypePartiallyFilled:
		if isFilledClose {
			side := futures.SideTypeBuy
			if order.Side == futures.SideTypeBuy {
				side = futures.SideTypeSell
			}
			_, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return b.bc.CloseOrder(ctx, order.ExecutedQuantity, side, order.Symbol)
			})
			if err != nil {
				b.l.Printf("error closing order: %v", err)
				b.handleError(err, userID, t.ID)
				return
			}
			// update trade state
			if changeState {
				b.c.UpdateTradeClosed(t.ID)
			}
			//message user
			msg := fmt.Sprintf("Order closed successfully.\n%s\n\nOrder ID: %d\nTrade ID: %d", strings.ToUpper(action), orderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
		}
		return
	case futures.OrderStatusTypeCanceled:
		return
	case futures.OrderStatusTypeExpired, futures.OrderStatusTypeRejected:
		return
	default:
		b.l.Printf("unknown order status received: %v", order.Status)

	}

}
