package bot

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

func (b *Bot) WsHandler(event *futures.WsUserDataEvent) {
	b.handleOrderTradeUpdate(context.Background(), event.OrderTradeUpdate)
}

func (b *Bot) WsErrHandler(err error) {
	b.l.Printf("WebSocket error: %v", err)
}

func (b *Bot) handleOrderTradeUpdate(ctx context.Context, f futures.WsOrderTradeUpdate) {
	switch f.Status {
	case futures.OrderStatusTypeCanceled:
		b.l.Println("Order was canceled.")
		// b.HandleCanceled(f)
	case futures.OrderStatusTypeFilled:
		b.l.Println("Order filled successfully.")
		// sleep a little bit to make sure the store is updated for early filled orders
		// TODO maybe change the logic in future for better handling
		time.Sleep(time.Second)
		go b.handleFilledExchange(ctx, f)
	case futures.OrderStatusTypeRejected:
		b.l.Println("Order was rejected.")
	case futures.OrderStatusTypeNew:
		b.l.Println("New order received.")
	case futures.OrderStatusTypeExpired:
		b.l.Println("Order has expired.")
	case futures.OrderStatusTypePartiallyFilled:
		b.l.Println("Order partially filled.")
	default:
		b.l.Println("Unknown order status received.")
	}
}
