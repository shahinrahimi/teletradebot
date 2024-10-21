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
	case types.ExchangeBinance:
		ex = b.bc
	case types.ExchangeBitmex:
		ex = b.mc
	default:
		msg := fmt.Sprintf("Unknown account: %s", t.Account)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
		return nil
	}

	b.DbgChan <- fmt.Sprintf("Closing trade: %d", t.ID)
	go func() {
		orderID := t.OrderID
		if orderID == "" {
			b.DbgChan <- fmt.Sprintf("Close main order skipped, orderID is empty, TradeID: %d", t.ID)
			return
		}
		b.DbgChan <- fmt.Sprintf("Closing main order of trade: %d", t.ID)
		b.handleCloseExchange(ctx, &t, i, types.OrderTitleMain, types.ExecutionCloseMainOrder, userID, orderID, ex, true, true)
	}()
	go func() {
		orderID := t.SLOrderID
		if orderID == "" {
			b.DbgChan <- fmt.Sprintf("Close stop-loss order skipped, orderID is empty, TradeID: %d", t.ID)
			return
		}
		b.DbgChan <- fmt.Sprintf("Closing stop-loss order of trade: %d", t.ID)
		b.handleCloseExchange(ctx, &t, i, types.OrderTitleStopLoss, types.ExecutionNone, userID, orderID, ex, false, false)
	}()
	go func() {
		orderID := t.TPOrderID
		if orderID == "" {
			b.DbgChan <- fmt.Sprintf("Close take-profit order skipped, orderID is empty, TradeID: %d", t.ID)
			return
		}
		b.DbgChan <- fmt.Sprintf("Closing take-profit order of trade: %d", t.ID)
		b.handleCloseExchange(ctx, &t, i, types.OrderTitleTakeProfit, types.ExecutionNone, userID, orderID, ex, false, false)
	}()
	go func() {
		orderID := t.ReverseOrderID
		if orderID == "" {
			b.DbgChan <- fmt.Sprintf("Close reverse-main order skipped, orderID is empty, TradeID: %d", t.ID)
			return
		}
		b.DbgChan <- fmt.Sprintf("Closing reverse-main order of trade: %d", t.ID)
		b.handleCloseExchange(ctx, &t, i, types.OrderTitleReverseMain, types.ExecutionCloseReverseMainOrder, userID, orderID, ex, true, true)
	}()
	go func() {
		orderID := t.ReverseSLOrderID
		if orderID == "" {
			b.DbgChan <- fmt.Sprintf("Close reverse-stop-loss order skipped, orderID is empty, TradeID: %d", t.ID)
			return
		}
		b.DbgChan <- fmt.Sprintf("Closing reverse-stop-loss order of trade: %d", t.ID)
		b.handleCloseExchange(ctx, &t, i, types.OrderTitleReverseStopLoss, types.ExecutionNone, userID, orderID, ex, false, false)
	}()
	go func() {
		orderID := t.ReverseTPOrderID
		if orderID == "" {
			b.DbgChan <- fmt.Sprintf("Close reverse-take-profit order skipped, orderID is empty, TradeID: %d", t.ID)
			return
		}
		b.DbgChan <- fmt.Sprintf("Closing reverse-take-profit order of trade: %d", t.ID)
		b.handleCloseExchange(ctx, &t, i, types.OrderTitleReverseTakeProfit, types.ExecutionNone, userID, orderID, ex, false, false)
	}()
	return nil
}
