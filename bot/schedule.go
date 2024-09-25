package bot

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) ScheduleOrderReplacement(ctx context.Context, i *models.Interpreter, t *models.Trade, ex exchange.Exchange) {
	b.c.SetInterpreter(i, t.ID)
	delay := i.CalculateExpiration()
	b.l.Printf("schedule order replacement: delay: %s, TradeID: %d", delay, t.ID)
	time.AfterFunc(delay, func() {
		oe := i.GetOrderExecution(types.GetOrderExecution, t.OrderID)
		res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return ex.GetOrder(ctx, oe)
		})
		if err != nil {
			b.l.Printf("error getting order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		status := utils.ExtractOrderStatus(res)

		switch status {
		case string(futures.OrderStatusTypeNew), string(futures.OrderStatusTypeExpired), swagger.OrderStatusTypeNew:
			b.l.Printf("Order not executed, attempting replacement, TradeID: %d", t.ID)
			oe := i.GetOrderExecution(types.CancelOrderExecution, t.OrderID)
			_, err := b.retryDenyNotFound(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return ex.CancelOrder(ctx, oe)
			})
			if err != nil {
				b.handleError(err, t.UserID, t.ID)
				return
			}
			time.Sleep(config.WaitForReplacement)
			// fetch new interpreter
			i, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return ex.FetchInterpreter(ctx, t)
			})
			if err != nil {
				b.l.Printf("error fetching interpreter: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			interpreter, ok := i.(*models.Interpreter)
			if !ok {
				b.l.Panicf("unexpected error happened in casting error to *models.Interpreter: %T", interpreter)
			}
			b.c.SetInterpreter(interpreter, t.ID)
			// place new order
			oe = interpreter.GetOrderExecution(types.StopPriceExecution, t.OrderID)
			res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return ex.PlaceStopOrder(ctx, oe)
			})
			if err != nil {
				b.l.Printf("error placing stop-order: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			// update trade
			b.c.UpdateTradeMainOrder(t.ID, res)
			// message user
			msg := b.getMessagePlacedOrder(types.OrderTitleMain, types.VerbReplaced, t.ID, res)
			b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
			// new schedule
			b.ScheduleOrderReplacement(ctx, interpreter, t, ex)
		default:
			b.l.Printf("Schedule order replacement canceled due to status: %s, TradeID: %d", status, t.ID)
		}
	})
}
