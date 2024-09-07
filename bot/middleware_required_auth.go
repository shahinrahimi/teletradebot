package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
)

func (b *Bot) RequiredAuth(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		var userID int64 = u.Message.From.ID
		for _, id := range config.UserIDs {
			if id == userID {
				next(u, ctx)
				return
			}
		}
		b.SendMessage(userID, "You are not allowed")
	}
}
