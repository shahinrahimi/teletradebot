package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleUnderDevelopment(u *tgbotapi.Update, ctx context.Context) error {
	//t := ctx.Value(models.KeyTrade{}).(models.Trade)
	b.SendMessage(u.Message.From.ID, "under development")
	return nil

}
