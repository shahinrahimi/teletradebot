package bot

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
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
	// fmt.Println("got an event")
	b.handleOrderTradeUpdate(event.OrderTradeUpdate)
}

func (b *Bot) errHandler(err error) {
	b.l.Printf("handling ws error: %v", err)
}

func (b *Bot) handleOrderTradeUpdate(f futures.WsOrderTradeUpdate) {
	// b.l.Printf("Start #######################(%s)########################", f.Status)
	// utils.PrintStructFields(f)
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
	// b.l.Printf("End #######################(%s)########################", f.Status)
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
	// check if canceled by bot schedule rather than the user manually cancel the order
	if t.State == types.STATE_REPLACING {
		po, err := b.bc.PrepareTradeForOrder(t)
		if err != nil {
			b.l.Printf("trade could not be executed, error in preparing state: %v", err)
			// update trade state to idle
			t.State = types.STATE_IDLE
			if err := b.s.UpdateTrade(t); err != nil {
				msg := fmt.Sprintf("An important error occurred. The trade with ID '%d' could not be updated, which might cause tracking issues. Order ID: %s", t.ID, t.OrderID)
				b.SendMessage(t.UserID, msg)
			}
			return
		}

		res, err := b.bc.PlacePreparedOrder(po)
		if err != nil {
			// TODO Handle binance api error
			// TODO send user a proper message base on binance API error
			b.l.Printf("error in placing trade: %v", err)
			return
		}

		b.SendMessage(t.UserID, "Order replaced with the new Order.")

		// schedule order cancellation (it will raise error if currently filled)
		// if cancel successfully it will change trade state to replacing
		go b.scheduleOrderCancellation(res.OrderID, res.Symbol, po.Expiration, t)

		// update trade state
		t.OrderID = strconv.FormatInt(res.OrderID, 10)
		t.State = types.STATE_PLACED
		t.UpdatedAt = time.Now().UTC()
		if err := b.s.UpdateTrade(t); err != nil {
			msg := fmt.Sprintf("An important error occurred. The trade with ID '%d' could not be updated, which might cause tracking issues. Order ID: %s", t.ID, t.OrderID)
			b.SendMessage(t.UserID, msg)
		}
		return
	}

}

func (b *Bot) scheduleOrderCancellation(orderID int64, symbol string, delay time.Duration, t *models.Trade) {
	time.AfterFunc(delay, func() {
		// check if trade has state of placed
		b.l.Println(t.State)
		if t.State == types.STATE_PLACED {
			// update trade state
			t.State = types.STATE_REPLACING
			if err := b.s.UpdateTrade(t); err != nil {
				b.l.Printf("important error occurred, failed to updated trade status fo replacing, the trade can not replaced: %v", err)
				return
			}
			// cancel order
			_, err := b.bc.CancelOrder(orderID, symbol)
			if err != nil {
				b.l.Printf("Failed to cancel order %d: %v", orderID, err)
				return
			}

		}

	})
}
