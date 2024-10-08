package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleHelp(u *tgbotapi.Update, ctx context.Context) error {
	userID := u.Message.From.ID
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: GetCommandHelp(),
	}
	return nil
}
