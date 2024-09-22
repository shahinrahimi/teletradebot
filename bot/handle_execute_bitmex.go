package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) HandleExecuteBitmex(update *tgbotapi.Update, ctx context.Context) {

	t, ok := ctx.Value(models.KeyTrade{}).(models.Trade)
	if !ok {
		b.l.Panic("error getting trade from context")
	}
	userID := update.Message.From.ID

	go func() {
		res, d, err := b.retry(config.MaxTries, config.WaitForNextTries, &t, func() (interface{}, interface{}, error) {
			return b.mc.PlaceTrade(ctx, &t)
		})
		if err != nil {
			b.l.Printf("error executing trade: %v", err)
			b.handleError(err, userID, t.ID)
			return
		}

		order, ok := res.(*swagger.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *bitmex.Order: %T", order)
			return
		}

		describer, ok := d.(*models.Describer)
		if !ok {
			b.l.Printf("unexpected error happened in casting interface to *models.Describer: %T", describer)
			return
		}
		b.c.SetDescriber(describer, t.ID)
		expiration := describer.CalculateExpiration()
		b.l.Printf("expiration: %v", utils.FriendlyDuration(expiration))
		// schedule for replacement
		go b.scheduleOrderReplacementBitmex(ctx, expiration, order.OrderID, &t)

		msg := fmt.Sprintf("Order placed successfully\n\nOrder ID: %s\nTrade ID: %d", order.OrderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}

		// update trade state
		b.c.UpdateTradePlaced(t.ID, order.OrderID)
	}()

}
