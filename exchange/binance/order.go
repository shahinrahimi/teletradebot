package binance

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
)

func (bc *BinanceClient) PlaceStopOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBinance)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", oe)
	}
	order := bc.client.NewCreateOrderService().
		Symbol(oeb.Symbol).
		Side(oeb.Side).
		Quantity(oeb.Quantity).
		StopPrice(oeb.StopPrice).
		Type(futures.OrderTypeStopMarket).
		WorkingType(futures.WorkingTypeMarkPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) PlaceTakeProfitOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBinance)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", oe)
	}
	order := bc.client.NewCreateOrderService().
		Symbol(oeb.Symbol).
		Side(oeb.Side).
		Quantity(oeb.Quantity).
		StopPrice(oeb.StopPrice).
		Type(futures.OrderTypeTakeProfitMarket).
		WorkingType(futures.WorkingTypeMarkPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) CancelOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBinance)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", oe)
	}
	order := bc.client.NewCancelOrderService().
		OrderID(oeb.OrderID).
		Symbol(oeb.Symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) GetOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBinance)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", oe)
	}
	order := bc.client.NewGetOrderService().
		OrderID(oeb.OrderID).
		Symbol(oeb.Symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) CloseOrder(ctx context.Context, oe interface{}) (interface{}, error) {
	oeb, ok := oe.(*models.OrderExecutionBinance)
	if !ok {
		bc.l.Panicf("unexpected order type: %T", oe)
	}
	order := bc.client.NewCreateOrderService().
		Symbol(oeb.Symbol).
		Side(oeb.Side).
		Type(futures.OrderTypeMarket).
		Quantity(oeb.Quantity).
		ReduceOnly(true)
	return order.Do(ctx)
}
