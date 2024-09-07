package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
)

func (b *Bot) HandleAdd(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if err := b.s.CreateTrade(&t); err != nil {
		b.l.Printf("error creating a new trade: %v", err)
		b.SendMessage(u.Message.From.ID, "Internal error while creating a new trade.")
		return err
	}
	b.SendMessage(u.Message.From.ID, "Trade created successfully!")
	return nil
}
