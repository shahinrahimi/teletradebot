package bot

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) handleFilledExchange(ctx context.Context, update interface{}) {
	var t *models.Trade
	var orderType types.OrderIDType
	var ex exchange.Exchange
	if u, ok := update.(futures.WsOrderTradeUpdate); ok {
		t, orderType = b.c.GetTradeByAnyOrderID(u.ID)
	} else if u, ok := update.(swagger.OrderData); ok {
		t, orderType = b.c.GetTradeByAnyOrderID(u.OrderID)
	} else {
		b.l.Panicf("unknown type of update: %T", update)
	}

	if types.OrderIDTypeNone == orderType {
		b.l.Printf("the orderID is not associate with any trade: %s", update)
		return
	}

	i, exist := b.c.GetInterpreter(t.ID)
	if !exist {
		b.l.Panicf("interpreter not exist for trade: %d", t.ID)
		return
	}

	if t.Account == string(types.ExchangeBinance) {
		ex = b.bc
	} else if t.Account == string(types.ExchangeBitmex) {
		ex = b.mc
	} else {
		b.l.Panicf("unknown account: %s", t.Account)
	}

	switch orderType {
	case types.OrderIDTypeMain:
		b.handleMainFilledExchange(ctx, t, i, ex)
	case types.OrderIDTypeTakeProfit:
		b.handleTakeProfitFilledExchange(ctx, t, i, ex)
	case types.OrderIDTypeStopLoss:
		b.handleStopLossFilledExchange(ctx, t, i, ex)
	case types.OrderIDTypeReverseMain:
		b.handleMainReversedFilledExchange(ctx, t, i, ex)
	case types.OrderIDTypeReverseTakeProfit:
		b.handleReverseTakeProfitFilledExchange(ctx, t, i, ex)
	case types.OrderIDTypeReverseStopLoss:
		b.handleReverseStopLossFilledExchange(ctx, t, i, ex)
	default:
		b.l.Printf("the orderID is not associate with any trade: %s", update)
	}
}

func (b *Bot) handleMainFilledExchange(ctx context.Context, t *models.Trade, i *models.Interpreter, ex exchange.Exchange) {
	// update trade state
	b.c.UpdateTradeFilled(t.ID)

	// message the user
	msg := b.getMessagePlacedOrder(types.OrderTitleMain, types.VerbFilled, t.ID, t.OrderID)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}

	// place stop-loss order
	go func() {
		oe := i.GetOrderExecution(types.StopLossExecution, "")
		res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return ex.PlaceStopOrder(ctx, oe)
		})
		if err != nil {
			b.l.Printf("error executing stop-loss order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		// update trade
		b.c.UpdateTradeSLOrder(t.ID, res)
		// message user
		msg := b.getMessagePlacedOrder(types.OrderTitleStopLoss, types.VerbPlaced, t.ID, res)
		b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
	}()

	// place take-profit order
	go func() {
		oe := i.GetOrderExecution(types.TakeProfitExecution, "")
		res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return ex.PlaceTakeProfitOrder(ctx, oe)
		})
		if err != nil {
			b.l.Printf("error executing take-profit order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		// update trade
		b.c.UpdateTradeTPOrder(t.ID, res)
		// message user
		msg := b.getMessagePlacedOrder(types.OrderTitleTakeProfit, types.VerbPlaced, t.ID, res)
		b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
	}()

	if t.ReverseMultiplier == 0 {
		return
	}
	// place reverse order
	go func() {
		oe := i.GetOrderExecution(types.ReverseStopPriceExecution, "")
		res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return b.bc.PlaceStopOrder(ctx, oe)
		})
		if err != nil {
			b.l.Printf("error executing reverse order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		// update trade
		b.c.UpdateTradeReverseOrder(t.ID, res)
		// message user
		msg := b.getMessagePlacedOrder(types.OrderTitleReverseMain, types.VerbPlaced, t.ID, res)
		b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
	}()
}

func (b *Bot) handleMainReversedFilledExchange(ctx context.Context, t *models.Trade, i *models.Interpreter, ex exchange.Exchange) {
	// update trade state
	b.c.UpdateTradeReverting(t.ID)

	// message the user
	msg := b.getMessagePlacedOrder(types.OrderTitleReverseMain, types.VerbFilled, t.ID, t.OrderID)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}

	// place stop-loss order
	go func() {
		oe := i.GetOrderExecution(types.ReverseStopPriceExecution, "")
		res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return ex.PlaceStopOrder(ctx, oe)
		})
		if err != nil {
			b.l.Printf("error executing stop-loss order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		// update trade
		b.c.UpdateTradeReverseSLOrder(t.ID, res)
		// message user
		msg := b.getMessagePlacedOrder(types.OrderTitleReverseStopLoss, types.VerbPlaced, t.ID, res)
		b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
	}()

	// place take-profit order
	go func() {
		oe := i.GetOrderExecution(types.TakeProfitExecution, "")
		res, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
			return ex.PlaceTakeProfitOrder(ctx, oe)
		})
		if err != nil {
			b.l.Printf("error executing take-profit order: %v", err)
			b.handleError(err, t.UserID, t.ID)
			return
		}
		// update trade
		b.c.UpdateTradeReverseTPOrder(t.ID, res)
		// message user
		msg := b.getMessagePlacedOrder(types.OrderTitleReverseTakeProfit, types.VerbPlaced, t.ID, res)
		b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
	}()
}

func (b *Bot) handleStopLossFilledExchange(ctx context.Context, t *models.Trade, i *models.Interpreter, ex exchange.Exchange) {
	if t.ReverseMultiplier == 0 {
		// update trade state
		b.c.UpdateTradeStopped(t.ID)
	}
	// message user
	msg := b.getMessageStopped(types.OrderTitleStopLoss, types.VerbFilled, t.ID)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}

	// cancel take-profit order
	go b.handleCancelExchange(ctx, t, i, types.OrderTitleTakeProfit, ex, t.TPOrderID)
}

func (b *Bot) handleTakeProfitFilledExchange(ctx context.Context, t *models.Trade, i *models.Interpreter, ex exchange.Exchange) {
	// update trade state
	b.c.UpdateTradeProfited(t.ID)
	// message user
	msg := b.getMessageStopped(types.OrderTitleTakeProfit, types.VerbFilled, t.ID)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}

	// cancel stop-loss order
	go b.handleCancelExchange(ctx, t, i, types.OrderTitleStopLoss, ex, t.SLOrderID)
}

func (b *Bot) handleReverseStopLossFilledExchange(ctx context.Context, t *models.Trade, i *models.Interpreter, ex exchange.Exchange) {
	// update trade state
	b.c.UpdateTradeStopped(t.ID)
	// message user
	msg := b.getMessageStopped(types.OrderTitleReverseStopLoss, types.VerbFilled, t.ID)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
	// cancel take-profit order
	go b.handleCancelExchange(ctx, t, i, types.OrderTitleReverseTakeProfit, ex, t.ReverseTPOrderID)
}

func (b *Bot) handleReverseTakeProfitFilledExchange(ctx context.Context, t *models.Trade, i *models.Interpreter, ex exchange.Exchange) {
	// update trade state
	b.c.UpdateTradeProfited(t.ID)
	// message user
	msg := b.getMessageStopped(types.OrderTitleReverseTakeProfit, types.VerbFilled, t.ID)
	b.MsgChan <- types.BotMessage{ChatID: t.UserID, MsgStr: msg}
	// cancel stop-loss order
	go b.handleCancelExchange(ctx, t, i, types.OrderTitleReverseStopLoss, ex, t.ReverseSLOrderID)
}
