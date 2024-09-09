package bot

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
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
		b.l.Println(f.ExecutionType) // for take profit is Trade
		b.l.Println(f.OriginalType)  // for take profit TAKE_PROFIT_MARKET
		b.l.Println(f.PositionSide)  // for take profit BOTH
		b.l.Println(f.ClientOrderID)
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
	orderID := utils.ConvertBinanceOrderID(f.ID)
	// check if order related to a trade
	t, err := b.s.GetTradeByOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			b.l.Printf("internal error for getting trade by OrderID")
		}
		utils.PrintStructFields(f)
		// probably the order created by another client
		return
	}
	if err := b.s.UpdateTradeFilled(t); err != nil {
		b.l.Printf("error updating a trade with new state: %s", types.STATE_FILLED)
		return
	}

	_, err1 := b.HandlePlaceStopLossOrder(t, &f)
	_, err2 := b.HandlePlaceTakeProfitOrder(t, &f)
	if err1 != nil {
		b.l.Printf("error placing stop loss : %v", err1)
		b.SendMessage(t.UserID, "could not place stop-loss order")
	} else {
		b.SendMessage(t.UserID, "stop-loss order placed successfully")
	}
	if err2 != nil {
		b.l.Printf("error placing take profit : %v", err2)
		b.SendMessage(t.UserID, "could not place take-profit order")
	} else {
		b.SendMessage(t.UserID, "take-profit order placed successfully")
	}

	// update trade
}

func (b *Bot) HandlePlaceStopLossOrder(t *models.Trade, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	psl, err := b.bc.PrepareStopLossOrder(context.Background(), t, f)
	if err != nil {
		return nil, fmt.Errorf("error placing stop-loss order in preparing stage: %v", err)
	}
	return b.bc.PlacePreparedStopLossOrder(context.Background(), psl)
}

func (b *Bot) HandlePlaceTakeProfitOrder(t *models.Trade, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	ptp, err := b.bc.PrepareTakeProfitOrder(context.Background(), t, f)
	if err != nil {
		return nil, fmt.Errorf("error placing take-profit order in preparing stage: %v", err)
	}
	return b.bc.PlacePreparedTakeProfitOrder(context.Background(), ptp)
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
				if err := b.s.UpdateTradeIdle(t); err != nil {
					b.l.Printf("error updating error: %v", err)
					return
				}
				return

			}
			NewOrderID := utils.ConvertBinanceOrderID(cp.OrderID)
			if err := b.s.UpdateTradePlaced(t, NewOrderID); err != nil {
				b.l.Printf("error updating trade: %v", err)
			}
			b.SendMessage(t.UserID, fmt.Sprintf("Order replaced successfully\nTrade ID: %d\nNewOrder ID: %d", t.ID, cp.OrderID))
			// schedule for replacement
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
