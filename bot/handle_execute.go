package bot

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange/binance"
	"github.com/shahinrahimi/teletradebot/exchange/bitmex"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) HandleExecute(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	if t.State != types.STATE_IDLE {
		msg := "The trade could not be executed as it has already been executed once."
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
		return nil
	}

	if t.Account == types.ACCOUNT_B {
		go func() {
			res, po, err := b.retry(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, interface{}, error) {
				return b.bc.PlaceTrade(ctx, &t)
			})
			if err != nil {
				b.l.Printf("error executing trade: %v", err)
				b.handleError(err, userID, t.ID)
				return
			}
			order, ok := res.(*futures.CreateOrderResponse)
			if !ok {
				b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", order)
				return
			}
			preparedOrder, ok := po.(*binance.PreparedOrder)
			if !ok {
				b.l.Printf("unexpected error happened in casting interface to futures.binance.PreparedOrder: %T", preparedOrder)
				return
			}
			orderID := utils.ConvertBinanceOrderID(order.OrderID)
			// schedule for replacement
			go b.scheduleOrderReplacement(ctx, preparedOrder.Expiration, order.OrderID, &t)

			msg := fmt.Sprintf("Order placed successfully\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: t.UserID,
				MsgStr: msg,
			}
			// update trade state
			if err := b.s.UpdateTradePlaced(&t, orderID); err != nil {
				b.l.Printf("error updating trade to DB: %v", err)
			}
		}()
	} else {

		go func() {
			res, po, err := b.retry(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, interface{}, error) {
				return b.mc.PlaceTrade(ctx, &t)
			})
			if err != nil {
				b.l.Printf("error executing trade: %v", err)
				b.handleError(err, userID, t.ID)
				return
			}

			order, ok := res.(*swagger.Order)
			if !ok {
				b.l.Printf("unexpected error happened in casting error to bitmex.Order: %T", order)
				return
			}

			preparedOrder, ok := po.(*bitmex.PreparedOrder)
			if !ok {
				b.l.Printf("unexpected error happened in casting interface to bitmex.PreparedOrder: %T", preparedOrder)
				return
			}

			msg := fmt.Sprintf("Order placed successfully\n\nOrder ID: %s\nTrade ID: %d", order.OrderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: t.UserID,
				MsgStr: msg,
			}

			// schedule for replacement
			go b.scheduleOrderReplacementBitmex(ctx, preparedOrder.Expiration, order.OrderID, &t)
			// update trade state
			if err := b.s.UpdateTradePlaced(&t, order.OrderID); err != nil {
				b.l.Printf("error updating trade to DB: %v", err)
			}
		}()

		//b.mc.GetLastClosedCandle(&t)
	}

	return nil
}
