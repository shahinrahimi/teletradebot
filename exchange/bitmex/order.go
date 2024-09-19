package bitmex

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	swagger "github.com/shahinrahimi/teletradebot/swagger"
)

func (mc *BitmexClient) PlacePreparedOrder(p *PreparedOrder) (*swagger.Order, error) {
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(p.Side),
		OrderQty: optional.NewFloat32(p.Quantity),
		OrdType:  optional.NewString("Stop"),
		StopPx:   optional.NewFloat64(p.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(mc.auth, p.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) PlaceOrder(ctx context.Context, p *PreparedOrder) (*swagger.Order, error) {
	ctx = mc.GetAuthContext(ctx)
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(p.Side),
		OrderQty: optional.NewFloat32(p.Quantity),
		OrdType:  optional.NewString(OrderTypeStop),
		StopPx:   optional.NewFloat64(p.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, p.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) GetOrder(ctx context.Context, symbol string, orderID string) (*swagger.Order, error) {
	ctx = mc.GetAuthContext(ctx)
	params := &swagger.OrderApiOrderGetOrdersOpts{
		Symbol: optional.NewString(symbol),
	}
	orders, _, err := mc.client.OrderApi.OrderGetOrders(ctx, params)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		if o.OrderID == orderID {
			return &o, nil
		}
	}
	return nil, fmt.Errorf("order not found")
}

func (mc *BitmexClient) CancelOrder(ctx context.Context, orderID string) (*swagger.Order, error) {
	ctx = mc.GetAuthContext(ctx)
	params := &swagger.OrderApiOrderCancelOpts{
		OrderID: optional.NewString(orderID),
	}
	orders, _, err := mc.client.OrderApi.OrderCancel(ctx, params)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		if o.OrderID == orderID {
			return &o, nil
		}
	}
	return nil, fmt.Errorf("order not found")
}
