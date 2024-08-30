package bot

import (
	"database/sql"
	"fmt"
	"strconv"

	"gihub.com/shahinrahimi/teletradebot/utils"
	"github.com/adshao/go-binance/v2/futures"
)

func (b *Bot) StartBinanceService() error {
	if err := b.bc.UpdateTickers(); err != nil {
		b.l.Printf("error updating tickers for binance : %v", err)
		return err
	}
	b.l.Printf("Total pairs found for binance: %d", len(b.bc.Symbols))

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

	doneC, _, err := futures.WsUserDataServe(b.bc.ListenKey, b.wsHandler, b.errHandler)
	if err != nil {
		b.l.Printf("error startUserDataStream: %v", err)
		return
	}

	b.l.Println("WebSocket connection established. Listening for events...")
	<-doneC
}

func (b *Bot) handleOrderTradeUpdate(f futures.WsOrderTradeUpdate) {
	utils.PrintStructFields(f)
	switch f.Status {
	case futures.OrderStatusTypeCanceled:
		b.l.Println("handle canceled")
	case futures.OrderStatusTypeFilled:
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

func (b *Bot) wsHandler(event *futures.WsUserDataEvent) {
	// handle order trade events
	b.handleOrderTradeUpdate(event.OrderTradeUpdate)
}

func (b *Bot) errHandler(err error) {
	b.l.Printf("handling ws error: %v", err)
}
