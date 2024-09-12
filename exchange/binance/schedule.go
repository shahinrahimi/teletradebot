package binance

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
)

func (bc *BinanceClient) scheduleOrderReplacement(ctx context.Context, delay time.Duration, orderId int64, t *models.Trade) {
	time.AfterFunc(delay, func() {
		order, err := bc.getOrder(ctx, orderId, t.Symbol)
		if err != nil {
			bc.l.Printf("Error retrieving order: %v", err)
			return
		}
		//b.l.Printf("Order ID: %d | Trade ID: %d | Order Status: %s", orderId, t.ID, order.Status)
		if order.Status == futures.OrderStatusTypeFilled {
			return
		}
		// order not executed
		if order.Status == futures.OrderStatusTypeNew {
			// cancel order
			if _, err := bc.cancelOrder(ctx, orderId, t.Symbol); err != nil {
				bc.l.Printf("Error canceling order: %v", err)
				bc.handleError(err, t.UserID)
				return
			}
			// update trade state to cancelled
			if err := bc.s.UpdateTradeCancelled(t); err != nil {
				bc.l.Printf("Error updating trade to CANCELED state: %v", err)
				return
			}
			// sleep a second to make sure the kline data is updated
			// TODO maybe need to change the logic
			time.Sleep(time.Second)
			// execute will tries for order placement and at the end will call scheduleOrderReplacement
			bc.ExecuteTrade(ctx, t, true)
		}
	})
}
