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
	if i, exist := b.c.GetInterpreter(t.ID); exist {
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: i.Describe(true),
		}
		return nil
	} else {

	}
	switch t.Account {
	case types.ExchangeBinance:
		go b.handleDescribeExchange(ctx, &t, userID, b.bc)
	case types.ExchangeBitmex:
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
