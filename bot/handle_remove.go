package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleRemove(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	if err := b.s.DeleteTrade(t.ID); err != nil {
		b.l.Printf("error deleting a trade: %v", err)
		return err
	}
	msg := fmt.Sprintf("The trade has been successfully removed.\n\nTrade ID: %d", t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: msg,
	}
	return nil
}
