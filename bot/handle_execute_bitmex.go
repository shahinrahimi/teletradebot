package bot

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) HandleExecuteBitmex(ctx context.Context, t *models.Trade, userID int64, ex exchange.Exchange) {

	i, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return ex.FetchInterpreter(ctx, t)
	})
	interpreter, ok := i.(*models.Interpreter)
	if !ok {
		b.l.Panicf("unexpected error happened in casting error to *models.Interpreter: %T", interpreter)
	}
	oe := interpreter.GetOrderExecutionBitmex(types.StopPriceExecution, t.OrderID)
	res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return ex.PlaceStopOrder(ctx, oe)
	})
	if err != nil {
		b.l.Printf("error placing stop-order: %v", err)
		b.handleError(err, userID, t.ID)
		return
	}
	order, ok := res.(*futures.CreateOrderResponse)
	if !ok {
		b.l.Panicf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", order)
	}

	orderID := utils.ConvertBinanceOrderID(order.OrderID)
	// schedule for replacement
	go b.scheduleOrderReplacementBitmex(ctx, interpreter, t, ex)

	msg := fmt.Sprintf("Order placed successfully\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}
	// update trade state
	b.c.UpdateTradeMainOrder(t.ID, orderID)

}
