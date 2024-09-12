package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleUnderDevelopment(u *tgbotapi.Update, ctx context.Context) error {
	//t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: "under development",
	}
	return nil

}
