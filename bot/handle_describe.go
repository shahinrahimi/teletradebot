package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
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
	case types.ACCOUNT_B:
		go b.handleFetchInterpreter(ctx, &t, userID, b.bc)
	case types.ACCOUNT_M:
		go b.handleFetchInterpreter(ctx, &t, userID, b.mc)
	default:
		msg := fmt.Sprintf("Unknown account: %s", t.Account)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
	}
	return nil
}

func (b *Bot) handleFetchInterpreter(ctx context.Context, t *models.Trade, userID int64, ex exchange.Exchange) {

	i, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return ex.FetchInterpreter(ctx, t)
	})
	if err != nil {
		b.l.Printf("error fetching interpreter: %v", err)
		b.handleError(err, userID, t.ID)
		return
	}
	interpreter, ok := i.(*models.Interpreter)
	if !ok {
		b.l.Panicf("unexpected error happened in casting error to *models.Interpreter: %T", interpreter)
	}
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: interpreter.Describe(),
	}
}
