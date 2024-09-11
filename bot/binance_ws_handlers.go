package bot

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/exchange/binance"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) handleOrderTradeUpdate(f futures.WsOrderTradeUpdate) {
	switch f.Status {
	case futures.OrderStatusTypeCanceled:
		b.l.Println("Order was canceled.")
		// b.HandleCanceled(f)
	case futures.OrderStatusTypeFilled:
		b.l.Println("Order filled successfully.")
		go b.handleFilled(f)
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

func (b *Bot) handleFilled(f futures.WsOrderTradeUpdate) {
	orderID := utils.ConvertBinanceOrderID(f.ID)
	var t *models.Trade
	var err error
	// sleep a little bit to make sure the store is updated for early filled orders
	// TODO maybe change the logic in future for better handling
	time.Sleep(time.Second)
	// check if order related to a trade
	t, err = b.s.GetTradeByOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			b.l.Panic("Internal error while retrieving trade:", err)
		}
	}
	// handle new filled order
	if t != nil {
		if err := b.s.UpdateTradeFilled(t); err != nil {
			b.l.Printf("Error updating trade state to FILLED: %v", types.STATE_FILLED)
			return
		}
		msg := fmt.Sprintf("Order filled successfully.\n\nTrade ID: %d", t.ID)
		b.SendMessage(t.UserID, msg)
		b.handleNewFilled(t, &f)
		return
	}
	// check if order is for stop loss
	t, err = b.s.GetTradeBySLOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			b.l.Panic("Internal error while retrieving stop-loss trade:", err)
		}
	}
	if t != nil {
		if err := b.s.UpdateTradeStopped(t); err != nil {
			b.l.Printf("Error updating trade state to STOPPED: %v", types.STATE_STOPPED)
			return
		}
		b.handleSLFilled(t, &f)
		return
	}
	// check if order is for take profit
	t, err = b.s.GetTradeByTPOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			b.l.Panic("Internal error while retrieving take-profit trade:", err)
		}
	}
	if t != nil {
		if err := b.s.UpdateTradeProfited(t); err != nil {
			b.l.Printf("Error updating trade state to PROFITED: %v", types.STATE_PROFITED)
			return
		}
		b.handleTPFilled(t, &f)
		return
	}

	// update trade
}
func (b *Bot) handleNewFilled(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	td, err := b.bc.GetTradeDescriber(context.Background(), t)
	if err != nil {
		b.l.Printf("error getting trade describer")
	}
	osl, err1 := b.handlePlaceStopLossOrder(t, f)
	otp, err2 := b.handlePlaceTakeProfitOrder(t, f)
	if err1 != nil {
		b.l.Printf("Error placing stop-loss order: %v", err1)
		msg := fmt.Sprintf("Failed to place stop-loss order.\n\nTrade ID: %d", t.ID)
		b.SendMessage(t.UserID, msg)
	} else {
		msg := fmt.Sprintf("Stop-loss order placed successfully.\n\nTrade ID: %d", t.ID)
		b.SendMessage(t.UserID, msg)
	}
	if err2 != nil {
		b.l.Printf("Error placing take-profit order: %v", err2)
		msg := fmt.Sprintf("Failed to place take-profit order.\n\nTrade ID: %d", t.ID)
		b.SendMessage(t.UserID, msg)
	} else {
		msg := fmt.Sprintf("Take-profit order placed successfully.\n\nTrade ID: %d", t.ID)
		b.SendMessage(t.UserID, msg)
	}

	if td != nil {
		binance.TradeDescribers[t.ID] = &binance.TradeDescriber{
			From:  td.From,
			Till:  td.Close,
			Open:  td.Open,
			Close: td.Close,
			High:  td.High,
			Low:   td.Low,
			TP:    otp.StopPrice,
			SL:    otp.StopPrice,
			SP:    f.StopPrice,
		}
	}

	slOrderID := utils.ConvertBinanceOrderID(osl.OrderID)
	tpOrderID := utils.ConvertBinanceOrderID(otp.OrderID)

	b.s.UpdateTradeSLandTP(t, slOrderID, tpOrderID)
}

func (b *Bot) handleSLFilled(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.TPOrderID)
	if err != nil {
		b.l.Panicf("Error converting take-profit order ID: %v", err)
		return
	}
	msg := fmt.Sprintf("🛑 Stop-loss order executed successfully.\n\nTrade ID: %d", t.ID)
	b.SendMessage(t.UserID, msg)
	if _, err := b.bc.CancelOrder(context.Background(), orderID, f.Symbol); err != nil {
		b.l.Printf("Error canceling take-profit order.")
		return
	}
	msg = fmt.Sprintf("Take-profit order has been canceled.\n\nTrade ID: %d", t.ID)
	b.SendMessage(t.UserID, msg)
}

func (b *Bot) handleTPFilled(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.SLOrderID)
	if err != nil {
		b.l.Panicf("Error converting stop-loss order ID: %v", err)
		return
	}
	msg := fmt.Sprintf("✅ Take-profit order executed successfully.\n\nTrade ID: %d", t.ID)
	b.SendMessage(t.UserID, msg)
	if _, err := b.bc.CancelOrder(context.Background(), orderID, f.Symbol); err != nil {
		b.l.Printf("Error canceling stop-loss order.: %v", err)
		return
	}
	msg = fmt.Sprintf("Stop-loss order has been canceled.\n\nTrade ID: %d", t.ID)
	b.SendMessage(t.UserID, msg)
}

func (b *Bot) handlePlaceStopLossOrder(t *models.Trade, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	psl, err := b.bc.PrepareStopLossOrder(context.Background(), t, f)
	if err != nil {
		return nil, fmt.Errorf("error during stop-loss order preparation: %v", err)
	}
	return b.bc.PlacePreparedStopLossOrder(context.Background(), psl)
}

func (b *Bot) handlePlaceTakeProfitOrder(t *models.Trade, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	ptp, err := b.bc.PrepareTakeProfitOrder(context.Background(), t, f)
	if err != nil {
		return nil, fmt.Errorf("error during take-profit order preparation: %v", err)
	}
	return b.bc.PlacePreparedTakeProfitOrder(context.Background(), ptp)
}

func (b *Bot) scheduleOrderReplacement(ctx context.Context, delay time.Duration, orderId int64, t *models.Trade) {
	time.AfterFunc(delay, func() {
		order, err := b.bc.GetOrder(ctx, orderId, t.Symbol)
		if err != nil {
			b.l.Printf("Error retrieving order: %v", err)
			return
		}
		//b.l.Printf("Order ID: %d | Trade ID: %d | Order Status: %s", orderId, t.ID, order.Status)
		if order.Status == futures.OrderStatusTypeFilled {
			return
		}
		// order not executed
		if order.Status == futures.OrderStatusTypeNew {
			// cancel order
			if _, err := b.bc.CancelOrder(ctx, orderId, t.Symbol); err != nil {
				b.l.Printf("Error canceling order: %v", err)
				b.handleAPIError(err, t.UserID)
				return
			}
			// update trade state to cancelled
			if err := b.s.UpdateTradeCancelled(t); err != nil {
				b.l.Printf("Error updating trade to CANCELED state: %v", err)
				return
			}
			// sleep a second to make sure the kline data is updated
			// TODO maybe need to change the logic
			time.Sleep(time.Second)
			// prepare new order
			p, err := b.bc.PrepareOrder(ctx, t)
			if err != nil {
				b.l.Printf("Error preparing new order: %v", err)
				return
			}
			// place new order
			cp, err := b.bc.PlacePreparedOrder(ctx, p)
			if err != nil {
				b.handleAPIError(err, t.UserID)
				return
			}
			NewOrderID := utils.ConvertBinanceOrderID(cp.OrderID)
			// update trade to placed
			if err := b.s.UpdateTradePlaced(t, NewOrderID); err != nil {
				b.l.Printf("Error updating trade to PLACED state: %v", err)
			}
			b.SendMessage(t.UserID, fmt.Sprintf("Order replaced successfully\n\nNewOrder ID: %d\nTrade ID: %d", cp.OrderID, t.ID))
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
