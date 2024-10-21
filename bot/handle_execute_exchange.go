package bot

import (
	"context"
	"fmt"

	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) handleExecuteExchange(ctx context.Context, t *models.Trade, userID int64, ex exchange.Exchange) {
	b.DbgChan <- fmt.Sprintf("Placing stop-order for trade: %d", t.ID)
	i, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return ex.FetchInterpreter(ctx, t)
	})
	if err != nil {
		b.handleError(err, userID, t.ID)
		return
	}
	interpreter, ok := i.(*models.Interpreter)
	if !ok {
		b.l.Panicf("unexpected error happened in casting res to *models.Interpreter: %T", interpreter)
	}
	oe := interpreter.GetOrderExecution(types.ExecutionEntryMainOrder, t.OrderID)
	b.DbgChan <- fmt.Sprintf("Placing stop-order with orderExecution: %v", oe)
	res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return ex.PlaceStopOrder(ctx, oe)
	})
	if err != nil {
		b.DbgChan <- fmt.Sprintf("error placing stop-order: %v", err)
		b.handleError(err, userID, t.ID)
		return
	}
	b.DbgChan <- fmt.Sprintf("Placed stop-order with result: %v", res)
	// update trade state
	b.c.UpdateTradeMainOrder(t.ID, res)
	// schedule for replacement
	go b.ScheduleOrderReplacement(ctx, interpreter, t.ID, ex)
	msg := b.getMessagePlacedOrder(types.OrderTitleMain, types.VerbPlaced, t.ID, res)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
}
