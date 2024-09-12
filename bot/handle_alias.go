package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
)

func (b *Bot) HandleAlias(u *tgbotapi.Update, ctx context.Context) error {
	userID := u.Message.From.ID
	var msg string = "aliases: \n"
	for key, value := range config.Shortcuts {
		msg = msg + "'" + key + "' => " + value + "\n"
	}
	b.MsgChan <- BotMessage{
		ChatID: userID,
		MsgStr: msg,
	}
	return nil
}
