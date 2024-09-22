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
	t, orderIDType := b.c.GetTradeByAnyOrderID(orderID)
	b.l.Printf("orderID: %s, orderIDType: %s", orderID, orderIDType)
	switch orderIDType {
	case types.OrderIDTypeMain:
		b.handleNewFilled(ctx, t, f)
	case types.OrderIDTypeTakeProfit:
		b.handleTPFilled(ctx, t, f)
	case types.OrderIDTypeStopLoss:
		b.handleSLFilled(ctx, t, f)
	default:
		b.l.Printf("the orderID is not associate with any trade: %s", orderID)
	}
}

func (b *Bot) handleNewFilled(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {

	// update trade state
	b.c.UpdateTradeFilled(t.ID)

	// message the user
	msg := fmt.Sprintf("Order filled successfully.\n\nTrade ID: %d", t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}

	d, exist := b.c.GetDescriber(t.ID)
	if !exist {
		b.l.Printf("describer not exist for trade: %d", t.ID)
		return
	}

	// place stop-loss order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.bc.PlaceTradeSLOrder(ctx, t, d, &f)
		})
		if err != nil {
			b.l.Printf("error executing stop-loss order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*futures.CreateOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", order)
			return
		}
		orderID := utils.ConvertBinanceOrderID(order.OrderID)
		// update trade
		b.c.UpdateTradeSLOrder(t.ID, orderID)

		// message the user
		msg := fmt.Sprintf("Stop-loss order placed successfully.\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()

	// place take-profit order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.bc.PlaceTradeTPOrder(ctx, t, d, &f)
		})
		if err != nil {
			b.l.Printf("error executing take-profit order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*futures.CreateOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", order)
			return
		}
		orderID := utils.ConvertBinanceOrderID(order.OrderID)
		// update trade
		b.c.UpdateTradeTPOrder(t.ID, orderID)

		// message the user
		msg := fmt.Sprintf("Take-profit order placed successfully.\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()
}

func (b *Bot) handleSLFilled(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {
	// update trade state
	b.c.UpdateTradeStopped(t.ID)

	// message the user
	msg := fmt.Sprintf("🛑 Stop-loss order executed successfully.\n\nPnL: %s\nTrade ID: %d", f.RealizedPnL, t.ID)
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
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*futures.CancelOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *futures.CancelOrderResponse: %T", order)
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
	b.c.UpdateTradeProfited(t.ID)

	// message the user
	msg := fmt.Sprintf("✅ Take-profit order executed successfully.\n\nPnL: %s\nTrade ID: %d", f.RealizedPnL, t.ID)
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
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*futures.CancelOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *futures.CancelOrderResponse: %T", order)
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

func (b *Bot) handleExpired(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {
	orderID := utils.ConvertBinanceOrderID(f.ID)
	t, orderIDType := b.c.GetTradeByAnyOrderID(orderID)
	switch orderIDType {
	case types.OrderIDTypeMain:
		b.l.Printf("the orderID is expired with trade ID: %s", orderID)
		//b.handleNewFilled(ctx, t, f)
	case types.OrderIDTypeTakeProfit:
		b.l.Printf("the take-profit orderID is expired with trade ID: %s", orderID)
		//b.handleTPFilled(ctx, t, f)
	case types.OrderIDTypeStopLoss:
		b.l.Printf("the stop-loss orderID is expired with trade ID: %s", orderID)
		//b.handleSLFilled(ctx, t, f)
	default:
		b.l.Printf("the orderID expired is not associate with any trade: %s", orderID)
	}
}
