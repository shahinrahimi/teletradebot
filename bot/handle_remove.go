package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
)

func (b *Bot) HandleRemove(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if err := b.s.DeleteTrade(t.ID); err != nil {
		b.l.Printf("error deleting a trade: %v", err)
		return err
	}
	msg := fmt.Sprintf("The trade has been successfully removed.\n\nTrade ID: %d", t.ID)
	b.SendMessage(u.Message.From.ID, msg)
	return nil
}
