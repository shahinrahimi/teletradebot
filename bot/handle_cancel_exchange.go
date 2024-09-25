package bot

import (
	"context"

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
	//oe := i.GetOrderExecution(types.GetOrderExecution,orderIDStr)
	oe := i.GetOrderExecution(types.CancelOrderExecution, orderIDStr)
	_, err := b.retryDenyNotFound(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return ex.CancelOrder(ctx, oe)
	})
	if err != nil {
		b.l.Printf("error canceling order: %v", err)
		b.handleError(err, t.UserID, t.ID)
		return
	}
	// message the user
	msg := b.getMessagePlacedOrder(orderTitle, types.VerbCanceled, t.ID, orderIDStr)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}

}
