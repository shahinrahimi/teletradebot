package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleRemove(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State != types.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "The trade could not be removed to it has not state of Idle")
		return nil
	}
	if err := b.s.DeleteTrade(t.ID); err != nil {
		return err
	}
	b.SendMessage(u.Message.From.ID, "Trade removed successfully!")
	return nil
}
