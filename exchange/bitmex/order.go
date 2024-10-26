package bitmex

import (
	"context"
	"fmt"
	"time"

	"github.com/antihax/optional"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
)

func (mc *BitmexClient) PlaceStopOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBitmex)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", oe)
	}
	ctx = mc.getAuthContext(ctx)
	mc.l.Printf("placing stop order: %+v", oeb.Quantity)
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(string(oeb.Side)),
		OrderQty: optional.NewFloat32(float32(oeb.Quantity)),
		OrdType:  optional.NewString(swagger.OrderTypeStop),
		StopPx:   optional.NewFloat64(oeb.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, oeb.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) PlaceTakeProfitOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBitmex)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", oe)
	}
	ctx = mc.getAuthContext(ctx)
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(string(oeb.Side)),
		OrderQty: optional.NewFloat32(float32(oeb.Quantity)),
		OrdType:  optional.NewString(swagger.OrderTypeMarketIfTouched),
		StopPx:   optional.NewFloat64(oeb.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, oeb.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) CancelOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBitmex)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", oe)
	}
	ctx = mc.getAuthContext(ctx)
	params := &swagger.OrderApiOrderCancelOpts{
		OrderID: optional.NewString(oeb.OrderID),
	}
	// TODO why order cancel from swagger returns a array of orders?
	orders, _, err := mc.client.OrderApi.OrderCancel(ctx, params)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		if o.OrderID == oeb.OrderID {
			return &o, nil
		}
	}

	return nil, fmt.Errorf("order not found")
}

func (mc *BitmexClient) GetOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBitmex)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", oe)
	}
	ctx = mc.getAuthContext(ctx)
	// endTime := time.Now().UTC().Format("2006-01-02 15:04")
	// filter := fmt.Sprintf(`{"endTime": "%s"}`, endTime)
	params := &swagger.OrderApiOrderGetOrdersOpts{
		Symbol:  optional.NewString(oeb.Symbol),
		Reverse: optional.NewBool(true),
		EndTime: optional.NewTime(time.Now().UTC()),
		// Filter:  optional.NewString(filter),
	}
	orders, _, err := mc.client.OrderApi.OrderGetOrders(ctx, params)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		if o.OrderID == oeb.OrderID {
			return &o, nil
		}
	}
	return nil, fmt.Errorf("order with id %s not found, total orders: %d", oeb.OrderID, len(orders))
}

func (mc *BitmexClient) CloseOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBitmex)
	if !ok {
		mc.l.Panicf("unexpected order type: %T", oe)
	}
	ctx = mc.getAuthContext(ctx)
	params := &swagger.OrderApiOrderNewOpts{
		// Symbol:  optional.NewString(symbol),
		ExecInst: optional.NewString("Close"),
		//OrderID: optional.NewString(orderID),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, oeb.Symbol, params)
	if err != nil {
		return nil, err
	}
	return &order, err
}
