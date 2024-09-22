package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleDescribeBitmex(u *tgbotapi.Update, ctx context.Context) {

	t, ok := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	if !ok {
		b.l.Panic("error getting trade from context")
	}
	go func() {
		d, err := b.mc.FetchDescriber(ctx, &t)
		if err != nil {
			b.l.Printf("error fetching describer: %v", err)
			return
		}
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: d.ToString(&t),
		}
	}()

}
