package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleInfo(u *tgbotapi.Update, ctx context.Context) error {
	userID := u.Message.From.ID
	username := u.Message.From.UserName
	msg := fmt.Sprintf("UserID:\t%d\nUsername:\t%s", userID, username)
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: msg,
	}
	return nil
}
