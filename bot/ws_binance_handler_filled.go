package bot

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

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
	case types.OrderIDTypeReverseMain:
		b.handleReverseFilled(ctx, t, f)
	case types.OrderIDTypeReverseTakeProfit:
		b.handleReverseTPFilled(ctx, t, f)
	case types.OrderIDTypeReverseStopLoss:
		b.handleReverseSLFilled(ctx, t, f)
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
			return b.bc.PlaceTradeSLOrder(ctx, d, &f)
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
			return b.bc.PlaceTradeTPOrder(ctx, d, &f)
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

	// place reverse order
	go func() {
		if !b.bc.ReverseEnabled {
			return
		}
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.bc.PlaceTradeReverseMainOrder(ctx, d, &f)
		})
		if err != nil {
			b.l.Printf("error executing reverse order: %v", err)
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
		b.c.UpdateTradeReverseOrder(t.ID, orderID)

		// message the user
		msg := fmt.Sprintf("Reverse order placed successfully.\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()

}

func (b *Bot) handleSLFilled(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {
	// update trade state
	if !b.bc.ReverseEnabled {
		b.c.UpdateTradeStopped(t.ID)
	}

	// message the user
	msg := fmt.Sprintf("🛑 Stop-loss order executed successfully.\n\nPnL: %s\nTrade ID: %d", f.RealizedPnL, t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}

	// cancel take-profit order
	go b.handleCloseBinanceOrder(ctx, t, t.UserID, "take-profit", t.TPOrderID, false, false)
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
	go b.handleCloseBinanceOrder(ctx, t, t.UserID, "stop-loss", t.SLOrderID, false, false)
}

func (b *Bot) handleReverseFilled(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {

	// update trade state
	b.c.UpdateTradeReverting(t.ID)

	d, exist := b.c.GetDescriber(t.ID)
	if !exist {
		b.l.Printf("describer not exist for trade: %d", t.ID)
		return
	}

	// place reverse stop-loss order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.bc.PlaceTradeReverseStopLossOrder(ctx, d, &f)
		})
		if err != nil {
			b.l.Printf("error executing reverse stop-loss order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*futures.CreateOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *futures.CreateOrderResponse: %T", order)
			return
		}
		orderID := utils.ConvertBinanceOrderID(order.OrderID)
		// update trade state
		b.c.UpdateTradeReverseSLOrder(t.ID, orderID)
		// message the user
		msg := fmt.Sprintf("Reverse stop-loss order placed successfully.\n\nOrderID: %s\nTrade ID: %d", orderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()

	// place reverse take-profit order
	go func() {
		res, err := b.retry2(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.bc.PlaceTradeReverseTakeProfitOrder(ctx, d, &f)
		})
		if err != nil {
			b.l.Printf("error executing reverse take-profit order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		order, ok := res.(*futures.CreateOrderResponse)
		if !ok {
			b.l.Printf("unexpected error happened in casting error to *futures.CreateOrderResponse: %T", order)
			return
		}
		orderID := utils.ConvertBinanceOrderID(order.OrderID)
		// update trade state
		b.c.UpdateTradeReverseTPOrder(t.ID, orderID)
		// message the user
		msg := fmt.Sprintf("Reverse take-profit order placed successfully.\n\nOrderID: %s\nTrade ID: %d", orderID, t.ID)
		b.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
	}()
}

func (b *Bot) handleReverseSLFilled(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {

	// update trade state
	b.c.UpdateTradeStopped(t.ID)

	// message the user
	msg := fmt.Sprintf("🛑 Stop-loss order executed successfully.\n\nPnL: %s\nTrade ID: %d", f.RealizedPnL, t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}

	// cancel take-profit order
	b.handleCloseBinanceOrder(ctx, t, t.UserID, "take-profit", t.ReverseTPOrderID, false, false)
}

func (b *Bot) handleReverseTPFilled(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {

	// update trade state
	b.c.UpdateTradeProfited(t.ID)

	// message the user
	msg := fmt.Sprintf("✅ Take-profit order executed successfully.\n\nPnL: %s\nTrade ID: %d", f.RealizedPnL, t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}

	// cancel stop-loss order
	go b.handleCloseBinanceOrder(ctx, t, t.UserID, "stop-loss", t.ReverseSLOrderID, false, false)
}
