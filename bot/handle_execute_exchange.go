package bot

import (
	"context"

	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) handleExecuteExchange(ctx context.Context, t *models.Trade, userID int64, ex exchange.Exchange) {
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
	res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return ex.PlaceStopOrder(ctx, oe)
	})
	if err != nil {
		b.l.Printf("error placing stop-order: %v", err)
		b.handleError(err, userID, t.ID)
		return
	}
	b.c.UpdateTradeMainOrder(t.ID, res)
	// schedule for replacement
	go b.ScheduleOrderReplacement(ctx, interpreter, t.ID, ex)
	msg := b.getMessagePlacedOrder(types.OrderTitleMain, types.VerbPlaced, t.ID, res)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
	// update trade state
}
