package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
)

func (b *Bot) HandleCloseBitmex(u *tgbotapi.Update, ctx context.Context) {

	t, ok := ctx.Value(models.KeyTrade{}).(models.Trade)
	if !ok {
		b.l.Panic("error getting trade from context")
	}
	userID := u.Message.From.ID

	go b.handleCloseBitmexOrder(ctx, &t, userID, "main", t.OrderID, true, true)
	go b.handleCloseBitmexOrder(ctx, &t, userID, "take-profit", t.TPOrderID, false, false)
	go b.handleCloseBitmexOrder(ctx, &t, userID, "stop-loss", t.SLOrderID, false, false)
	go b.handleCloseBitmexOrder(ctx, &t, userID, "reverse", t.ReverseOrderID, true, true)
	go b.handleCloseBitmexOrder(ctx, &t, userID, "reverse-take-profit", t.ReverseTPOrderID, false, false)
	go b.handleCloseBitmexOrder(ctx, &t, userID, "reverse-stop-loss", t.ReverseSLOrderID, false, false)
}
