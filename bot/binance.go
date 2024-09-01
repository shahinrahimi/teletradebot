package bot

import (
	"database/sql"
	"fmt"
	"strconv"

	"gihub.com/shahinrahimi/teletradebot/types"
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
	b.l.Printf("Start #######################(%s)########################", f.Status)
	utils.PrintStructFields(f)
	switch f.Status {
	case futures.OrderStatusTypeCanceled:
		b.l.Println("handle canceled")
		b.HandleCanceled(f)
	case futures.OrderStatusTypeFilled:
		b.l.Println("handle filled")
		b.HandleFilled(f)
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
	b.l.Printf("End #######################(%s)########################", f.Status)
}

func (b *Bot) HandleFilled(f futures.WsOrderTradeUpdate) {
	orderID := strconv.FormatInt(f.ID, 10)
	// check if order related to a trade
	t, err := b.s.GetTradeByOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			b.l.Printf("internal error for getting trade by OrderID")
		}
		// probably the order created by another client
		return
	}
	// update trade status
	t.State = types.STATE_FILLED
	if err := b.s.UpdateTrade(t); err != nil {
		b.l.Printf("error updating a trade with new state: %s", types.STATE_FILLED)
		b.SendMessage(t.UserID, "Could not update the trade with the new state.")
		return
	}

	res1, res2, err1, err2 := b.bc.TryPlaceStopLossAndTakeProfitTrade(t, &f)
	if err1 != nil || err2 != nil {
		b.l.Printf("error placing sl and tp orders: %v %v", err1, err2)
		return
	}
	msg := fmt.Sprintf("SL and TP orders placed successfully. Order IDs: %d, %d", res1.OrderID, res2.OrderID)
	b.SendMessage(t.UserID, msg)
}

func (b *Bot) HandleCanceled(f futures.WsOrderTradeUpdate) {
	orderID := strconv.FormatInt(f.ID, 10)
	// check if order related to a trade
	t, err := b.s.GetTradeByOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			b.l.Printf("internal error for getting trade by OrderID")
		}
		// probably the order created by another client
		return
	}
	res, err := b.bc.TryPlaceOrderForTrade(t)
	if err != nil {
		b.l.Printf("the order canceled but can't to replace a another order :%v", err)
		b.SendMessage(t.UserID, "Could not replace the order.")
		return
	}
	t.OrderID = strconv.FormatInt(res.OrderID, 10)
	if err := b.s.UpdateTrade(t); err != nil {
		b.l.Printf("error updating a trade with new order_id: %v", err)
		b.SendMessage(t.UserID, "Could not update the trade with the new order ID.")
		return
	}
	b.SendMessage(t.UserID, "Order replaced with the new order ID.")

}
