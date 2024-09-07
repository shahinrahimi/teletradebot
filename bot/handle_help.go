package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleHelp(u *tgbotapi.Update, ctx context.Context) error {
	var userID int64 = u.Message.From.ID
	b.SendMessage(userID, GetCommandHelp())
	return nil
}
