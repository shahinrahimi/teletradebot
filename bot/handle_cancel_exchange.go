package bot

import (
	"context"
	"fmt"

	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) handleCancelExchange(
	ctx context.Context,
	t *models.Trade,
	i *models.Interpreter,
	orderTitle types.OrderTitleType,
	ex exchange.Exchange,
	orderIDStr string,
) {
	if orderIDStr == "" {
		b.DbgChan <- fmt.Sprintf("Cancel order skipped, orderIDStr is empty, TradeID: %d", t.ID)
		return
	}
	//oe := i.GetOrderExecution(types.GetOrderExecution,orderIDStr)
	oe := i.GetOrderExecution(types.ExecutionCancelOrder, orderIDStr)
	_, err := b.retryDenyNotFound(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		b.DbgChan <- fmt.Sprintf("Canceling order: %s, TradeID: %d", orderIDStr, t.ID)
		return ex.CancelOrder(ctx, oe)
	})
	if err != nil {
		b.DbgChan <- fmt.Sprintf("error canceling order: %v, TradeID: %d", err, t.ID)
		b.handleError(err, t.UserID, t.ID)
		return
	}
	// message the user
	b.DbgChan <- fmt.Sprintf("Order canceled: %s, TradeID: %d", orderIDStr, t.ID)
	msg := b.getMessagePlacedOrder(orderTitle, types.VerbCanceled, t.ID, orderIDStr)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}

}
