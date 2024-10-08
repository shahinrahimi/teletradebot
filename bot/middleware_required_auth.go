package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) RequiredAuth(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		userID := u.Message.From.ID
		username := u.Message.From.UserName
		for _, id := range config.UserIDs {
			if id == userID {
				next(u, ctx)
				return
			}
		}
		msg := fmt.Sprintf("You are not allowed\n\n Username: %s UserID: %d", username, userID)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
	}
}
