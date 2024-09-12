package binance

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
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
