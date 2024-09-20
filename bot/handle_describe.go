package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleDescribe(u *tgbotapi.Update, ctx context.Context) error {
	t, ok := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	if !ok {
		b.l.Panic("error getting trade from context")
	}
	if d, exist := b.c.GetDescriber(t.ID); exist {
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: d.ToString(&t),
		}
		return nil
	}
	switch t.Account {
	case types.ACCOUNT_B:
		d, err := b.bc.FetchDescriber(ctx, &t)
		if err != nil {
			b.l.Printf("error fetching describer")
			return err
		}
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: d.ToString(&t),
		}
		return nil
	case types.ACCOUNT_M:
		d, err := b.mc.FetchDescriber(ctx, &t)
		if err != nil {
			b.l.Printf("error fetching describer")
			return err
		}
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: d.ToString(&t),
		}
		return nil
	}
	return nil
}
