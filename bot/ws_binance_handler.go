package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) wsHandler(event *futures.WsUserDataEvent) {
	b.handleOrderTradeUpdate(context.Background(), event.OrderTradeUpdate)
}

func (b *Bot) errHandler(err error) {
	b.l.Printf("WebSocket error: %v", err)
}

func (b *Bot) GetAnyTradeAssociateWithOrderID(orderID string) (trade *models.Trade, isOrderID bool, isTPOrderID bool, isSLOrderID bool, err error) {
	ts, err := b.s.GetTrades()
	if err != nil {
		return nil, false, false, false, err
	}
	for _, t := range ts {
		switch {
		case t.OrderID == orderID:
			return t, true, false, false, nil
		case t.TPOrderID == orderID:
			return t, false, true, false, nil
		case t.SLOrderID == orderID:
			return t, false, false, true, nil
		default:
			continue
		}
	}
	return nil, false, false, false, nil
}

func (b *Bot) handleOrderTradeUpdate(ctx context.Context, f futures.WsOrderTradeUpdate) {
	switch f.Status {
	case futures.OrderStatusTypeCanceled:
		b.l.Println("Order was canceled.")
		// b.HandleCanceled(f)
	case futures.OrderStatusTypeFilled:
		b.l.Println("Order filled successfully.")
		// sleep a little bit to make sure the store is updated for early filled orders
		// TODO maybe change the logic in future for better handling
		time.Sleep(time.Second)
		go b.handleFilled(ctx, f)
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

func (b *Bot) handleFilled(ctx context.Context, f futures.WsOrderTradeUpdate) {
	orderID := utils.ConvertBinanceOrderID(f.ID)
	t, isOrderID, isTPOrderID, isSLOrderID, err := b.GetAnyTradeAssociateWithOrderID(orderID)
	if err != nil {
		b.l.Printf("error getting associate trade with ID: %v", err)
	}
	switch {
	case isOrderID:
		b.handleNewFilled(ctx, t, f)
	case isTPOrderID:
		b.handleTPFilled(ctx, t, f)
	case isSLOrderID:
		b.handleSLFilled(ctx, t, f)
	default:
		b.l.Printf("the orderID is not associate with any trade: %v", err)
	}
}

func (b *Bot) handleNewFilled(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {
	d, err := b.bc.FetchDescriber(context.Background(), t)
	if err != nil {
		b.l.Printf("error fetching the describer %v", err)
	} else {
		models.SetDescriber(d, t.ID)
	}

	// update trade state
	if err := b.s.UpdateTradeFilled(t); err != nil {
		b.l.Panic("Internal error while updating trade:", err)
		return
	}
	// message the user
	msg := fmt.Sprintf("Order filled successfully.\n\nTrade ID: %d", t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}

	// place stop-loss order
	go func() {
		res, _, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, interface{}, error) {
			return b.bc.PlaceTradeSLOrder(ctx, t, &f)
		})
		if err != nil {
			b.l.Printf("error executing stop-loss order: %v", err)
			b.handleError(err, t.UserID)
			return
		}
		orderResponse, ok := res.(*futures.CreateOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", orderResponse)
			return
		}
		// message the user
		msg := fmt.Sprintf("Stop-loss order placed successfully.\n\nOrder ID: %d\nTrade ID: %d", orderResponse.OrderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()

	// place take-profit order
	go func() {
		res, _, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, interface{}, error) {
			return b.bc.PlaceTradeTPOrder(ctx, t, &f)
		})
		if err != nil {
			b.l.Printf("error executing take-profit order: %v", err)
			b.handleError(err, t.UserID)
			return
		}
		orderResponse, ok := res.(*futures.CreateOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", orderResponse)
			return
		}
		// message the user
		msg := fmt.Sprintf("Take-profit order placed successfully.\n\nOrder ID: %d\nTrade ID: %d", orderResponse.OrderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()
}

func (b *Bot) handleSLFilled(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {
	// update trade state
	if err := b.s.UpdateTradeStopped(t); err != nil {
		b.l.Panic("Internal error while updating trade:", err)
		return
	}
	// message the user
	msg := fmt.Sprintf("🛑 Stop-loss order executed successfully.\n\nTrade ID: %d", t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}

	// cancel take-profit order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.TPOrderID)
			if err != nil {
				return nil, err
			}
			return b.bc.CancelOrder(ctx, orderID, f.Symbol)
		})
		if err != nil {
			b.l.Printf("error cancelling take-profit order: %v", err)
			b.handleError(err, t.UserID)
			return
		}
		orderResponse, ok := res.(*futures.CancelOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to futures.CancelOrderResponse: %T", orderResponse)
			return
		}
		// message the user
		msg = fmt.Sprintf("Take-profit order has been canceled.\n\nTrade ID: %d", t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()
}

func (b *Bot) handleTPFilled(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {
	// update trade state
	if err := b.s.UpdateTradeProfited(t); err != nil {
		b.l.Panic("Internal error while updating trade:", err)
		return
	}
	// message the user
	msg := fmt.Sprintf("✅ Take-profit order executed successfully.\n\nTrade ID: %d", t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}

	// cancel stop-loss order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.SLOrderID)
			if err != nil {
				return nil, err
			}
			return b.bc.CancelOrder(ctx, orderID, f.Symbol)
		})
		if err != nil {
			b.l.Printf("error cancelling take-profit order: %v", err)
			b.handleError(err, t.UserID)
			return
		}
		orderResponse, ok := res.(*futures.CancelOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to futures.CancelOrderResponse: %T", orderResponse)
			return
		}
		// message the user
		msg = fmt.Sprintf("Stop-loss order has been canceled.\n\nTrade ID: %d", t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()
}