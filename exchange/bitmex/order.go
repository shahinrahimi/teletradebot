package bitmex

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	swagger "github.com/shahinrahimi/teletradebot/swagger"
)

func (mc *BitmexClient) PlaceOrder(ctx context.Context, po *PreparedOrder) (*swagger.Order, error) {
	ctx = mc.GetAuthContext(ctx)
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(string(po.Side)),
		OrderQty: optional.NewFloat32(float32(po.Quantity)),
		OrdType:  optional.NewString(OrderTypeStop),
		StopPx:   optional.NewFloat64(po.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, po.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) PlaceSLOrder(ctx context.Context, po *PreparedOrder) (*swagger.Order, error) {
	ctx = mc.GetAuthContext(ctx)
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(string(po.Side)),
		OrderQty: optional.NewFloat32(float32(po.Quantity)),
		OrdType:  optional.NewString(OrderTypeStop),
		StopPx:   optional.NewFloat64(po.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, po.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) PlaceTPOrder(ctx context.Context, po *PreparedOrder) (*swagger.Order, error) {
	ctx = mc.GetAuthContext(ctx)
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(string(po.Side)),
		OrderQty: optional.NewFloat32(float32(po.Quantity)),
		OrdType:  optional.NewString(OrderTypeStop),
		StopPx:   optional.NewFloat64(po.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, po.Symbol, params)
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
