package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) Logger(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		if u.Message.Command() == "" {
			b.l.Printf("Received message: %s", u.Message.Text)
		} else {
			b.l.Printf("Received command: %s\t args: %s", u.Message.Command(), u.Message.CommandArguments())
		}
		next(u, ctx)
	}
}
