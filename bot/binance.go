package bot

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
	"github.com/adshao/go-binance/v2/futures"
)

func (b *Bot) StartBinanceService(ctx context.Context) {
	go b.startUserDataStream(ctx)
}

func (b *Bot) startUserDataStream(ctx context.Context) {
	futures.UseTestnet = b.bc.UseTestnet
	listenKey, err := b.bc.GetListenKey(ctx)
	if err != nil {
		b.l.Printf("error getting listen key: %v", err)
		return
	}

	doneC, stopC, err := futures.WsUserDataServe(listenKey, b.wsHandler, b.errHandler)
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
		// b.HandleCanceled(f)
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

	if psl, err := b.bc.PrepareStopLossOrder(context.Background(), t, &f); err == nil {
		_, err = b.bc.PlacePreparedTakeProfitOrder(context.Background(), psl)
		if err != nil {
			b.SendMessage(t.UserID, "could not place stop-loss order")
		} else {
			b.SendMessage(t.UserID, "stop-loss order placed successfully")
		}
	} else {
		b.l.Printf("error placing stop-loss order in preparing stage: %v", err)
	}
	if ptp, err := b.bc.PrepareTakeProfitOrder(context.Background(), t, &f); err == nil {
		_, err = b.bc.PlacePreparedTakeProfitOrder(context.Background(), ptp)
		if err != nil {
			b.SendMessage(t.UserID, "could not place take-profit order")
		} else {
			b.SendMessage(t.UserID, "take-profit order placed successfully")
		}
	} else {
		b.l.Printf("error placing take-profit order in preparing stage: %v", err)
	}
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
		fmt.Println("replacing")
		// first prepare new order
		po, err := b.bc.PrepareOrder(context.Background(), t)
		if err != nil {
			b.l.Printf("error replacing order in preparing stage: %v", err)
			t.State = types.STATE_IDLE
			t.UpdatedAt = time.Now().UTC()
			t.OrderID = ""
			if err := b.s.UpdateTrade(t); err != nil {
				b.l.Printf("error updating trade state: %v", err)
			}
			return
		}
		// place order
		res, err := b.bc.PlacePreparedOrder(context.Background(), po)
		if err != nil {
			b.l.Printf("error replacing order in placing stage: %v", err)
			t.State = types.STATE_IDLE
			t.UpdatedAt = time.Now().UTC()
			t.OrderID = ""
			if err := b.s.UpdateTrade(t); err != nil {
				b.l.Printf("error updating trade state: %v", err)
			}
			return
		}

		b.SendMessage(t.UserID, "Order replaced with the new Order.")

		// schedule order cancellation (it will raise error if currently filled)
		// if cancel successfully it will change trade state to replacing
		go b.scheduleOrderCancellation(res.OrderID, po.Expiration, t)

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

func (b *Bot) scheduleOrderCancellation(orderID int64, delay time.Duration, t *models.Trade) {
	time.AfterFunc(delay, func() {
		// check if trade has state of placed
		if t.State == types.STATE_PLACED {
			// update trade state
			t.State = types.STATE_REPLACING
			t.UpdatedAt = time.Now().UTC()
			t.OrderID = ""
			if err := b.s.UpdateTrade(t); err != nil {
				b.l.Printf("important error occurred, failed to updated trade status fo replacing, the trade can not replaced: %v", err)
				return
			}
			// cancel order
			_, err := b.bc.CancelOrder(context.Background(), orderID, t.Symbol)
			if err != nil {
				b.l.Printf("Failed to cancel order %d: %v", orderID, err)
				return
			}

		}

	})
}

func (b *Bot) scheduleOrderReplacement(ctx context.Context, delay time.Duration, orderId int64, t *models.Trade) {
	time.AfterFunc(delay, func() {
		order, err := b.bc.GetOrder(ctx, orderId, t.Symbol)
		if err != nil {
			b.l.Printf("error getting order: %v", err)
			return
		}
		b.l.Printf("order_id: %d trade_id %d order_status %s", orderId, t.ID, order.Status)
		if order.Status == futures.OrderStatusTypeFilled {
			return
		}
		// order not executed
		if order.Status == futures.OrderStatusTypeNew {
			// cancel order
			if _, err := b.bc.CancelOrder(ctx, orderId, t.Symbol); err != nil {
				b.l.Printf("error canceling order: %v", err)
				b.handleAPIError(err, t.UserID)
				return
			}
			// prepare new order
			p, err := b.bc.PrepareOrder(ctx, t)
			if err != nil {
				b.l.Printf("error preparing order: %v", err)
				return
			}
			// place new order
			cp, err := b.bc.PlacePreparedOrder(ctx, p)
			if err != nil {
				b.handleAPIError(err, t.UserID)
				t.OrderID = ""
				t.UpdatedAt = time.Now().UTC()
				t.State = types.STATE_IDLE
				if err := b.s.UpdateTrade(t); err != nil {
					b.l.Printf("error updating error: %v", err)
					return
				}
				return
			}
			b.SendMessage(t.UserID, fmt.Sprintf("Order replaced successfully\nTrade ID: %d\n NewOrder ID: %d", t.ID, cp.OrderID))
			// update trade
			t.OrderID = strconv.FormatInt(cp.OrderID, 10)
			t.UpdatedAt = time.Now().UTC()
			if err := b.s.UpdateTrade(t); err != nil {
				b.l.Printf("error updating error: %v", err)
				return
			}
			go b.scheduleOrderReplacement(ctx, p.Expiration, cp.OrderID, t)
		}
	})
}

func (b *Bot) ScanningTrades(ctx context.Context) {
	ts, err := b.s.GetTrades()
	if err != nil {
		return
	}
	for _, t := range ts {
		if t.State != types.STATE_IDLE {
			if t.Account == types.ACCOUNT_B {
				orderID, err := strconv.ParseInt(t.OrderID, 10, 64)
				if err != nil {
					continue
				}
				order, err := b.bc.GetOrder(ctx, orderID, t.Symbol)
				if err != nil || order.Status == futures.OrderStatusTypeCanceled {
					t.State = types.STATE_IDLE
					t.OrderID = ""
					b.s.UpdateTrade(t)
				}

			}

		}

	}
}
