package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleAdd(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	id, err := b.s.CreateTrade(&t)
	if err != nil {
		b.l.Panicf("error creating a new trade: %v", err)
	}
	msg := fmt.Sprintf("Trade created successfully!\n\n Trade ID: %d", id)
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: msg,
	}
	return nil
}
