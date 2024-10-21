package bot

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) handleCloseExchange(
	ctx context.Context,
	t *models.Trade,
	i *models.Interpreter,
	orderTitle types.OrderTitleType,
	closeOrderType types.ExecutionType,
	userID int64,
	orderIDStr string,
	ex exchange.Exchange,
	changeState bool,
	isFilledClose bool,
) {
	if orderIDStr == "" {
		b.DbgChan <- fmt.Sprintf("Close order skipped, orderIDStr is empty, TradeID: %d", t.ID)
		return
	}
	oe := i.GetOrderExecution(types.ExecutionGetOrder, orderIDStr)

	// get order
	res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return ex.GetOrder(ctx, oe)
	})
	if err != nil {
		b.DbgChan <- fmt.Sprintf("error getting order: %v", err)
		b.handleError(err, userID, t.ID)
		return
	}
	status := utils.ExtractOrderStatus(res)

	b.DbgChan <- fmt.Sprintf("order status is: %s", status)

	switch status {
	case string(futures.OrderStatusTypeNew), swagger.OrderStatusTypeNew:
		oe := i.GetOrderExecution(types.ExecutionCancelOrder, orderIDStr)
		_, err := b.retryDenyNotFound(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return ex.CancelOrder(ctx, oe)
		})
		if err != nil {
			b.DbgChan <- fmt.Sprintf("error cancelling order: %v", err)
			b.handleError(err, userID, t.ID)
			return
		}

		// update trade state
		if changeState {
			b.c.UpdateTradeCanceled(t.ID)
		}

		// message user
		msg := b.getMessagePlacedOrder(orderTitle, types.VerbCanceled, t.ID, orderIDStr)
		b.MsgChan <- types.BotMessage{ChatID: userID, MsgStr: msg}
		return
	case string(futures.OrderStatusTypeFilled), string(futures.OrderStatusTypePartiallyFilled), swagger.OrderStatusTypeFilled:
		if isFilledClose {
			oe := i.GetOrderExecution(closeOrderType, orderIDStr)
			_, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
				return ex.CloseOrder(ctx, oe)
			})
			if err != nil {
				b.DbgChan <- fmt.Sprintf("error closing order: %v", err)
				b.handleError(err, userID, t.ID)
				return
			}
			// update trade state
			if changeState {
				b.c.UpdateTradeClosed(t.ID)
			}
			//message user
			msg := b.getMessagePlacedOrder(orderTitle, types.VerbClosed, t.ID, orderIDStr)
			b.MsgChan <- types.BotMessage{ChatID: userID, MsgStr: msg}
		}
		return
	case string(futures.OrderStatusTypeCanceled), swagger.OrderStatusTypeCanceled:
		return
	case string(futures.OrderStatusTypeExpired), string(futures.OrderStatusTypeRejected):
		return
	default:
		b.DbgChan <- fmt.Sprintf("unknown order status received: %v", status)

	}
}
