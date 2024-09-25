package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleClose(u *tgbotapi.Update, ctx context.Context) error {
	t, ok := ctx.Value(models.KeyTrade{}).(models.Trade)
	if !ok {
		b.l.Panic("error getting trade from context")
	}

	userID := u.Message.From.ID
	i, exist := b.c.GetInterpreter(t.ID)
	if !exist {
		return nil
	}
	var ex exchange.Exchange
	switch t.Account {
	case string(types.ExchangeBinance):
		ex = b.bc
	case string(types.ExchangeBitmex):
		ex = b.mc
	default:
		msg := fmt.Sprintf("Unknown account: %s", t.Account)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
		return nil
	}

	go b.handleCloseExchange(ctx, &t, i, types.OrderTitleMain, types.ClosePositionExecution, userID, t.OrderID, ex, true, true)
	go b.handleCloseExchange(ctx, &t, i, types.OrderTitleStopLoss, types.ExecutionNone, userID, t.OrderID, ex, false, false)
	go b.handleCloseExchange(ctx, &t, i, types.OrderTitleTakeProfit, types.ExecutionNone, userID, t.OrderID, ex, false, false)
	go b.handleCloseExchange(ctx, &t, i, types.OrderTitleReverseMain, types.ClosePositionReverseExecution, userID, t.OrderID, ex, true, true)
	go b.handleCloseExchange(ctx, &t, i, types.OrderTitleReverseStopLoss, types.ExecutionNone, userID, t.OrderID, ex, false, false)
	go b.handleCloseExchange(ctx, &t, i, types.OrderTitleReverseTakeProfit, types.ExecutionNone, userID, t.OrderID, ex, false, false)
	return nil
}
