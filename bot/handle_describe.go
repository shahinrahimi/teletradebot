package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleDescribe(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.Account == types.ACCOUNT_B {
		td, err := b.bc.GetTradeDescriber(ctx, &t)
		if err != nil {
			return err
		}
		b.SendMessage(u.Message.From.ID, td.ToTelegramString(&t))

	} else {
		return b.HandleUnderDevelopment(u, ctx)
	}
	return nil
}
