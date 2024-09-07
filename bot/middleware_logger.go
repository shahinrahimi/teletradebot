package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) Logger(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		var userID int64 = u.Message.From.ID
		if u.Message.Command() == "" {
			b.l.Printf("Received message (%d): %s", userID, u.Message.Text)
		} else {
			b.l.Printf("Received command (%d): %s\t args: %s", userID, u.Message.Command(), u.Message.CommandArguments())
		}
		next(u, ctx)
	}
}
