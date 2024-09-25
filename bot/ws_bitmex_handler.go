package bot

import (
	"context"
	"time"

	"github.com/shahinrahimi/teletradebot/swagger"
)

func (b *Bot) WsHandlerBitmex(ctx context.Context, od []swagger.OrderData) {
	for _, o := range od {
		b.l.Printf("Order got status: %s , orderType: %s", o.OrdStatus, o.OrdType)
		switch o.OrdStatus {
		case swagger.OrderStatusTypeNew:
			b.l.Println("New order received.")
		case swagger.OrderStatusTypeFilled:
			time.Sleep(time.Second)
			b.handleFilledBitmex(ctx, o)
		default:
			b.l.Println("Unknown order status received.")
		}
	}
}
