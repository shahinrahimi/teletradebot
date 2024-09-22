package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange/bitmex"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) handleOrderTradeUpdateBitmex(ctx context.Context, od []bitmex.OrderData) {
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

func (b *Bot) handleFilledBitmex(ctx context.Context, od bitmex.OrderData) {
	t, orderIDType := b.c.GetTradeByAnyOrderID(od.OrderID)
	switch orderIDType {
	case types.OrderIDTypeMain:
		b.handleNewFilledBitmex(ctx, t, od)
	case types.OrderIDTypeTakeProfit:
		b.handleTPFilledBitmex(ctx, t, od)
	case types.OrderIDTypeStopLoss:
		b.handleSLFilledBitmex(ctx, t, od)
	default:
		b.l.Printf("the orderID is not associate with any trade: %s", od.OrderID)
	}
}

func (b *Bot) handleNewFilledBitmex(ctx context.Context, t *models.Trade, od bitmex.OrderData) {
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
			return b.mc.PlaceTradeSLOrder(ctx, t, d, &od)
		})
		if err != nil {
			b.l.Printf("error executing stop-loss order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*swagger.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", order)
			return
		}
		orderID := order.OrderID
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
			return b.mc.PlaceTradeTPOrder(ctx, t, d, &od)
		})
		if err != nil {
			b.l.Printf("error executing take-profit order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*swagger.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to futures.CreateOrderResponse: %T", order)
			return
		}
		orderID := order.OrderID
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

func (b *Bot) handleSLFilledBitmex(ctx context.Context, t *models.Trade, od bitmex.OrderData) {
	// update trade state
	b.c.UpdateTradeStopped(t.ID)

	// message the user
	msg := fmt.Sprintf("🛑 Stop-loss order executed successfully.\n\nTrade ID: %d", t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}

	// cancel take-profit order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.mc.CancelOrder(ctx, t.TPOrderID)
		})
		if err != nil {
			b.l.Printf("error cancelling take-profit order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*swagger.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", order)
			return
		}
		// message the user
		msg = fmt.Sprintf("Take-profit order has been canceled.\n\nTrade ID: %d", t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()

	// place new stop-loss order

}

func (b *Bot) handleTPFilledBitmex(ctx context.Context, t *models.Trade, od bitmex.OrderData) {
	// update trade state
	b.c.UpdateTradeProfited(t.ID)

	// message the user
	msg := fmt.Sprintf("✅ Take-profit order executed successfully.\n\nTrade ID: %d", t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}

	// cancel stop-loss order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.mc.CancelOrder(ctx, t.SLOrderID)
		})
		if err != nil {
			b.l.Printf("error cancelling take-profit order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*swagger.Order)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *swagger.Order: %T", order)
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
