package binance

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) placeOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		StopPrice(p.StopPrice).
		Type(futures.OrderTypeStopMarket).
		WorkingType(futures.WorkingTypeMarkPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) placeTPOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		WorkingType(futures.WorkingTypeMarkPrice).
		Type(futures.OrderTypeTakeProfitMarket).
		StopPrice(p.StopPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) placeSLOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		WorkingType(futures.WorkingTypeMarkPrice).
		Type(futures.OrderTypeStopMarket).
		StopPrice(p.StopPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) cancelOrder(ctx context.Context, orderID int64, symbol string) (*futures.CancelOrderResponse, error) {
	order := bc.Client.NewCancelOrderService().
		OrderID(orderID).
		Symbol(symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) getOrder(ctx context.Context, orderID int64, symbol string) (*futures.Order, error) {
	order := bc.Client.NewGetOrderService().
		OrderID(orderID).
		Symbol(symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) getTradeLatestDescriber(ctx context.Context, t *models.Trade) (*types.TradeDescriber, error) {
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	if err != nil {
		return nil, err
	}
	sp, err := bc.calculateStopPrice(t, k, s)
	if err != nil {
		return nil, err
	}
	sl, err := bc.calculateStopLossPrice(t, k, s, sp)
	if err != nil {
		return nil, err
	}
	tp, err := bc.calculateTakeProfitPrice(t, k, s, sp)
	if err != nil {
		return nil, err
	}

	from := utils.ConvertTime(k.OpenTime)
	till := utils.ConvertTime(k.CloseTime).Add(time.Second)

	return &types.TradeDescriber{
		From:  from,
		Till:  till,
		Open:  k.Open,
		Close: k.Close,
		High:  k.High,
		Low:   k.Low,
		SP:    sp,
		TP:    tp,
		SL:    sl,
	}, nil
}

func (bc *BinanceClient) GetTradeDescriber(ctx context.Context, t *models.Trade) (*types.TradeDescriber, error) {
	td, exist := types.TradeDescribers[t.ID]
	if exist {
		return td, nil
	}
	return bc.getTradeLatestDescriber(ctx, t)
}
