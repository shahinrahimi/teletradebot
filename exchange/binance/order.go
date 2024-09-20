package binance

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
)

func (bc *BinanceClient) PlaceOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		StopPrice(p.StopPrice).
		Type(futures.OrderTypeStopMarket).
		WorkingType(futures.WorkingTypeMarkPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) PlaceTPOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		WorkingType(futures.WorkingTypeMarkPrice).
		Type(futures.OrderTypeTakeProfitMarket).
		StopPrice(p.StopPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) PlaceSLOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		WorkingType(futures.WorkingTypeMarkPrice).
		Type(futures.OrderTypeStopMarket).
		StopPrice(p.StopPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) CancelOrder(ctx context.Context, orderID int64, symbol string) (*futures.CancelOrderResponse, error) {
	order := bc.Client.NewCancelOrderService().
		OrderID(orderID).
		Symbol(symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) GetOrder(ctx context.Context, orderID int64, symbol string) (*futures.Order, error) {
	order := bc.Client.NewGetOrderService().
		OrderID(orderID).
		Symbol(symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) CloseOrder(ctx context.Context, quantity string, side futures.SideType, symbol string) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(futures.OrderTypeMarket).
		Quantity(quantity).
		ReduceOnly(true)

	return order.Do(ctx)
}
