package bot

import (
	"context"
	"time"

	"github.com/shahinrahimi/teletradebot/exchange/bitmex"
)

func (b *Bot) handleOrderTradeUpdateBitmex(ctx context.Context, od []OrderData) {
	for _, o := range od {
		b.l.Printf("Order got status: %s , orderType: %s", o.OrdStatus, o.OrdType)
		switch o.OrdStatus {
		case bitmex.OrderStatusTypeNew:
			b.l.Println("New order received.")
		case bitmex.OrderStatusTypeFilled:
			time.Sleep(time.Second)
			b.handleFilledBitmex(ctx, o)
		default:
			b.l.Println("Unknown order status received.")
		}
	}
}

func (b *Bot) handleFilledBitmex(ctx context.Context, o OrderData) {
	_, isOrderID, isTPOrderID, isSLOrderID, err := b.findTradeWithAnyOrderID(o.OrderID)
	if err != nil {
		b.l.Printf("error getting associate trade with ID: %v", err)
	}
	switch {
	case isOrderID:
		//b.handleNewFilledBitmex(ctx, t, o)
	case isTPOrderID:
		//b.handleTPFilledBitmex(ctx, t, o)
	case isSLOrderID:
		//b.handleSLFilledBitmex(ctx, t, o)
	default:
		b.l.Printf("the orderID is not associate with any trade: %v", err)
	}
}

// func (b *Bot) handleNewFilledBitmex(ctx context.Context, t *models.Trade, o OrderData) {
// 	// d, err := b.mc.FetchDescriber(context.Background(), t)
// 	// if err != nil {
// 	// 	b.l.Printf("error fetching the describer %v", err)
// 	// } else {
// 	// 	models.SetDescriber(d, t.ID)
// 	// }
// 	// update trade state
// 	if err := b.s.UpdateTradeFilled(t); err != nil {
// 		b.l.Panic("Internal error while updating trade:", err)
// 		return
// 	}
// 	// message the user
// 	msg := fmt.Sprintf("Order filled successfully.\n\nTrade ID: %d", t.ID)
// 	b.MsgChan <- types.BotMessage{
// 		ChatID: t.UserID,
// 		MsgStr: msg,
// 	}

// 	// place stop-loss order
// 	go func() {
// 		res, _, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, interface{}, error) {
// 			return b.bc.PlaceTradeSLOrder(ctx, t, &f)
// 		})
// 		if err != nil {
// 			b.l.Printf("error executing stop-loss order: %v", err)
// 			b.handleError(err, t.UserID, t.ID)
// 			return
// 		}
// 		orderResponse, ok := res.(*futures.CreateOrderResponse)
// 		if !ok {
// 			b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", orderResponse)
// 			return
// 		}
// 		orderID := utils.ConvertBinanceOrderID(orderResponse.OrderID)
// 		// update trade
// 		if err := b.s.UpdateTradeSLOrder(t, orderID); err != nil {
// 			b.l.Printf("error updating trade state: %v", err)
// 			return
// 		}
// 		// message the user
// 		msg := fmt.Sprintf("Stop-loss order placed successfully.\n\nOrder ID: %d\nTrade ID: %d", orderResponse.OrderID, t.ID)
// 		b.MsgChan <- types.BotMessage{
// 			ChatID: t.UserID,
// 			MsgStr: msg,
// 		}
// 	}()

// 	// place take-profit order
// 	go func() {
// 		res, _, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, interface{}, error) {
// 			return b.bc.PlaceTradeTPOrder(ctx, t, &f)
// 		})
// 		if err != nil {
// 			b.l.Printf("error executing take-profit order: %v", err)
// 			b.handleError(err, t.UserID, t.ID)
// 			return
// 		}
// 		orderResponse, ok := res.(*futures.CreateOrderResponse)
// 		if !ok {
// 			b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", orderResponse)
// 			return
// 		}
// 		orderID := utils.ConvertBinanceOrderID(orderResponse.OrderID)
// 		// update trade
// 		if err := b.s.UpdateTradeTPOrder(t, orderID); err != nil {
// 			b.l.Printf("error updating trade state: %v", err)
// 			return
// 		}
// 		// message the user
// 		msg := fmt.Sprintf("Take-profit order placed successfully.\n\nOrder ID: %d\nTrade ID: %d", orderResponse.OrderID, t.ID)
// 		b.MsgChan <- types.BotMessage{
// 			ChatID: t.UserID,
// 			MsgStr: msg,
// 		}
// 	}()
// }

// func (b *Bot) handleSLFilledBitmex(ctx context.Context, t *models.Trade, o OrderData) {
// 	// update trade state
// 	if err := b.s.UpdateTradeStopped(t); err != nil {
// 		b.l.Panic("Internal error while updating trade:", err)
// 		return
// 	}
// 	// message the user
// 	msg := fmt.Sprintf("🛑 Stop-loss order executed successfully.\n\nPnL: %s\nTrade ID: %d", f.RealizedPnL, t.ID)
// 	b.MsgChan <- types.BotMessage{
// 		ChatID: t.UserID,
// 		MsgStr: msg,
// 	}

// 	// cancel take-profit order
// 	go func() {
// 		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
// 			orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.TPOrderID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return b.bc.CancelOrder(ctx, orderID, f.Symbol)
// 		})
// 		if err != nil {
// 			b.l.Printf("error cancelling take-profit order: %v", err)
// 			b.handleError(err, t.UserID, t.ID)
// 			return
// 		}
// 		orderResponse, ok := res.(*futures.CancelOrderResponse)
// 		if !ok {
// 			b.l.Printf("unexpected error happened in casting error to futures.CancelOrderResponse: %T", orderResponse)
// 			return
// 		}
// 		// message the user
// 		msg = fmt.Sprintf("Take-profit order has been canceled.\n\nTrade ID: %d", t.ID)
// 		b.MsgChan <- types.BotMessage{
// 			ChatID: t.UserID,
// 			MsgStr: msg,
// 		}
// 	}()
// }

// func (b *Bot) handleTPFilledBitmex(ctx context.Context, t *models.Trade, o OrderData) {
// 	// update trade state
// 	if err := b.s.UpdateTradeProfited(t); err != nil {
// 		b.l.Panic("Internal error while updating trade:", err)
// 		return
// 	}
// 	// message the user
// 	msg := fmt.Sprintf("✅ Take-profit order executed successfully.\n\nPnL: %s\nTrade ID: %d", f.RealizedPnL, t.ID)
// 	b.MsgChan <- types.BotMessage{
// 		ChatID: t.UserID,
// 		MsgStr: msg,
// 	}

// 	// cancel stop-loss order
// 	go func() {
// 		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
// 			orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.SLOrderID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return b.bc.CancelOrder(ctx, orderID, f.Symbol)
// 		})
// 		if err != nil {
// 			b.l.Printf("error cancelling take-profit order: %v", err)
// 			b.handleError(err, t.UserID, t.ID)
// 			return
// 		}
// 		orderResponse, ok := res.(*futures.CancelOrderResponse)
// 		if !ok {
// 			b.l.Printf("unexpected error happened in casting error to futures.CancelOrderResponse: %T", orderResponse)
// 			return
// 		}
// 		// message the user
// 		msg = fmt.Sprintf("Stop-loss order has been canceled.\n\nTrade ID: %d", t.ID)
// 		b.MsgChan <- types.BotMessage{
// 			ChatID: t.UserID,
// 			MsgStr: msg,
// 		}
// 	}()
// }
