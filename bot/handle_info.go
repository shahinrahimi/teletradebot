package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleInfo(u *tgbotapi.Update, ctx context.Context) error {
	var userID int64 = u.Message.From.ID
	var username string = u.Message.From.UserName
	msg := fmt.Sprintf("UserID:\t%d\nUsername:\t%s", userID, username)
	b.SendMessage(userID, msg)
	return nil
}
