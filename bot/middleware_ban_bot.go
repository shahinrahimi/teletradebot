package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) BanBots(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		if u.Message.From.IsBot {
			return
		}
		next(u, ctx)
	}
}
