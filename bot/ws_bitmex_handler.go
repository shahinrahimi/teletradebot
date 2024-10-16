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
			go func() {
				time.Sleep(time.Second)
				b.handleFilledExchange(ctx, o)
			}()
		case swagger.OrderStatusTypeCanceled:
			go func() {
				time.Sleep(time.Second * 2)
				b.l.Println("Order was canceled.")
				b.l.Printf("Order status: %s , orderType: %s, orderRejectionReason: %s", o.OrdStatus, o.OrdType, o.OrdRejReason)
				b.l.Printf("Order: %+v", o)
				b.handleCanceledExchange(ctx, o)
			}()

		default:
			b.l.Printf("Unhandled order status received: %s", o.OrdStatus)
		}
	}
}
