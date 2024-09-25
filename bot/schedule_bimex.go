package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) scheduleOrderReplacementBitmex(ctx context.Context, i *models.Interpreter, t *models.Trade, ex exchange.Exchange) {
	b.c.SetInterpreter(i, t.ID)
	delay := i.CalculateExpiration()
	b.l.Printf("schedule order replacement: delay: %s, TradeID: %d", delay, t.ID)
	time.AfterFunc(delay, func() {
		oe := i.GetOrderExecutionBitmex(types.GetOrderExecution, t.OrderID)
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return ex.GetOrder(ctx, oe)
		})
		if err != nil {
			b.l.Printf("error getting order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := (res).(*swagger.Order)
		if !ok {
			b.l.Panicf("unexpected error happened in casting error to futures.Order: %T", order)
		}

		switch order.OrdStatus {
		case swagger.OrderStatusTypeNew:
			b.l.Printf("Order not executed, attempting replacement, TradeID: %d", t.ID)
			oe := i.GetOrderExecutionBitmex(types.CancelOrderExecution, t.OrderID)
			_, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return ex.CancelOrder(ctx, oe)
			})
			if err != nil {
				// TODO add handleError for cases that canceled and the orderId not found
				b.handleError(err, t.UserID, t.ID)
				b.c.UpdateTradeCanceled(t.ID)
				return
			}
			time.Sleep(config.WaitForReplacement)
			// fetch new interpreter
			i, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
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
			oe = interpreter.GetOrderExecutionBitmex(types.StopPriceExecution, t.OrderID)
			res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return ex.PlaceStopOrder(ctx, oe)
			})
			if err != nil {
				b.l.Printf("error placing stop-order: %v", err)
				b.handleError(err, t.UserID, t.ID)
				return
			}
			order, ok := (res).(*swagger.Order)
			if !ok {
				b.l.Panicf("unexpected error happened in casting error to futures.Order: %T", order)
			}
			// update trade
			b.c.UpdateTradeMainOrder(t.ID, order.OrderID)
			// message user
			msg := fmt.Sprintf("Order replaced successfully.\n\nNewOrder ID: %s\nTrade ID: %d", order.OrderID, t.ID)
			b.MsgChan <- types.BotMessage{
				ChatID: t.UserID,
				MsgStr: msg,
			}
			// new schedule
			b.scheduleOrderReplacementBitmex(ctx, interpreter, t, b.mc)
		default:
			b.l.Printf("Schedule order replacement canceled due to status: %s, TradeID: %d", order.OrdStatus, t.ID)
		}
	})
}
