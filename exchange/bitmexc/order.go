package bitmexc

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"github.com/shahinrahimi/teletradebot/swagger"
)

type PreparedOrder struct {
	OrderID   string // use for get and cancel the order
	Symbol    string
	Quantity  float64
	StopPrice float64
	Side      SideType
}

func (mc *BitmexClient) PlaceStopOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", po)
	}
	ctx = mc.getAuthContext(ctx)
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(string(p.Side)),
		OrderQty: optional.NewFloat32(float32(p.Quantity)),
		OrdType:  optional.NewString(OrderTypeStop),
		StopPx:   optional.NewFloat64(p.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, p.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) PlaceTakeProfitOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", po)
	}
	ctx = mc.getAuthContext(ctx)
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(string(p.Side)),
		OrderQty: optional.NewFloat32(float32(p.Quantity)),
		OrdType:  optional.NewString(OrderTypeMarketIfTouched),
		StopPx:   optional.NewFloat64(p.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, p.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) CancelOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", po)
	}
	ctx = mc.getAuthContext(ctx)
	params := &swagger.OrderApiOrderCancelOpts{
		OrderID: optional.NewString(p.OrderID),
	}
	// TODO why order cancel from swagger returns a array of orders?
	orders, _, err := mc.client.OrderApi.OrderCancel(ctx, params)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		if o.OrderID == p.OrderID {
			return &o, nil
		}
	}

	return nil, fmt.Errorf("order not found")
}

func (mc *BitmexClient) GetOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", po)
	}
	ctx = mc.getAuthContext(ctx)
	params := &swagger.OrderApiOrderGetOrdersOpts{
		Symbol: optional.NewString(p.Symbol),
	}
	orders, _, err := mc.client.OrderApi.OrderGetOrders(ctx, params)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		if o.OrderID == p.OrderID {
			return &o, nil
		}
	}
	return nil, fmt.Errorf("order not found")
}

func (mc *BitmexClient) CloseOrder(ctx context.Context, po interface{}) (interface{}, error) {
	p, ok := po.(*PreparedOrder)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", po)
	}
	ctx = mc.getAuthContext(ctx)
	params := &swagger.OrderApiOrderNewOpts{
		// Symbol:  optional.NewString(symbol),
		ExecInst: optional.NewString("Close"),
		//OrderID: optional.NewString(orderID),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, p.Symbol, params)
	if err != nil {
		return nil, err
	}
	return &order, err
}
