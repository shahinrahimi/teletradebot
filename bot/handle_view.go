package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
)

func (b *Bot) HandleView(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	b.MsgChan <- BotMessage{
		ChatID: userID,
		MsgStr: t.ToViewString(),
	}
	return nil
}
