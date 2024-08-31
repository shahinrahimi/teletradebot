package bot

import (
	"database/sql"
	"fmt"
	"strconv"

	"gihub.com/shahinrahimi/teletradebot/utils"
	"github.com/adshao/go-binance/v2/futures"
)

func (b *Bot) StartBinanceService() error {
	if err := b.bc.UpdateListenKey(); err != nil {
		b.l.Printf("error updating listenKey for binance : %v", err)
		return err
	}
	b.l.Printf("ListenKey acquired: %s", b.bc.ListenKey)
	go b.startUserDataStream()
	return nil
}

func (b *Bot) startUserDataStream() {
	futures.UseTestnet = b.bc.UseTestnet

	doneC, stopC, err := futures.WsUserDataServe(b.bc.ListenKey, b.wsHandler, b.errHandler)
	if err != nil {
		b.l.Printf("error startUserDataStream: %v", err)
		return
	}
	defer close(stopC)

	b.l.Println("WebSocket connection established. Listening for events...")
	<-doneC
}

func (b *Bot) wsHandler(event *futures.WsUserDataEvent) {
	// handle order trade events
	fmt.Println("got an event")
	b.handleOrderTradeUpdate(event.OrderTradeUpdate)
}

func (b *Bot) errHandler(err error) {
	b.l.Printf("handling ws error: %v", err)
}

func (b *Bot) handleOrderTradeUpdate(f futures.WsOrderTradeUpdate) {
	utils.PrintStructFields(f)
	switch f.Status {
	case futures.OrderStatusTypeCanceled:
		b.l.Println("handle canceled")
	case futures.OrderStatusTypeFilled:
		orderID := strconv.FormatInt(f.ID, 10)
		t, err := b.s.GetTradeByOrderID(orderID)
		if err != nil {
			if err != sql.ErrNoRows {
				b.l.Printf("internal error for getting trade by OrderID")
			}
			// probably the order created by another client
			return
		}
		res1, res2, err1, err2 := b.bc.TryPlaceStopLossAndTakeProfitTrade(t, &f)
		if err1 != nil || err2 != nil {
			b.l.Printf("error placing sl and tp orders: %v %v", err1, err2)
		}
		msg := fmt.Sprintf("the sl and tp order placed => order_id: %d order_id: %d", res1.OrderID, res2.OrderID)
		b.SendMessage(t.UserID, msg)
		// update trade to filled
		// place stopPrice order
		// place takeprofit order
		b.l.Println("handle filled")
	case futures.OrderStatusTypeRejected:
		b.l.Println("handle rejected")
	case futures.OrderStatusTypeNew:
		b.l.Println("handle new order")
	case futures.OrderStatusTypeExpired:
		b.l.Println("handle expiration")
	case futures.OrderStatusTypePartiallyFilled:
		b.l.Println("handle partially filled")
	default:
		b.l.Println("handle unknown")
	}
	orderID := strconv.FormatInt(f.ID, 10)
	t, err := b.s.GetTradeByOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			b.l.Printf("internal error for getting trade by OrderID")
		}
		// probably the order created by another client
		return
	}
	msg := fmt.Sprintf("the order status changed order_id: %s", orderID)
	b.SendMessage(t.UserID, msg)

}
