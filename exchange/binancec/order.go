package binancec

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
)

type PreparedOrder struct {
	OrderID   int64 // use for get and cancel the order
	Symbol    string
	Quantity  string
	StopPrice string
	Side      futures.SideType
}

func (bc *BinanceClient) PlaceStopOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", po)
	}
	order := bc.client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		StopPrice(p.StopPrice).
		Type(futures.OrderTypeStopMarket).
		WorkingType(futures.WorkingTypeMarkPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) PlaceTakeProfitOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", po)
	}
	order := bc.client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		StopPrice(p.StopPrice).
		Type(futures.OrderTypeTakeProfitMarket).
		WorkingType(futures.WorkingTypeMarkPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) CancelOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", po)
	}
	order := bc.client.NewCancelOrderService().
		OrderID(p.OrderID).
		Symbol(p.Symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) GetOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", po)
	}
	order := bc.client.NewGetOrderService().
		OrderID(p.OrderID).
		Symbol(p.Symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) CloseOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", po)
	}
	order := bc.client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Type(futures.OrderTypeMarket).
		Quantity(p.Quantity).
		ReduceOnly(true)
	return order.Do(ctx)
}
