package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
)

func (b *Bot) HandleAlias(u *tgbotapi.Update, ctx context.Context) error {
	var msg string = "aliases: \n"
	for key, value := range config.Shortcuts {
		msg = msg + "'" + key + "' => " + value + "\n"
	}
	b.SendMessage(u.Message.From.ID, msg)
	return nil
}
