package bot

import (
	"context"
	"fmt"

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
			MsgStr: d.ToString(),
		}
		return nil
	} else {

	}
	switch t.Account {
	case string(types.ExchangeBinance):
		go b.handleDescribeExchange(ctx, &t, userID, b.bc)
	case string(types.ExchangeBitmex):
		go b.handleDescribeExchange(ctx, &t, userID, b.mc)
	default:
		msg := fmt.Sprintf("Unknown account: %s", t.Account)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
	}
	return nil
}
