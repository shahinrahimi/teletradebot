package binance

import (
	"context"
	"fmt"
	"time"

	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
	"gihub.com/shahinrahimi/teletradebot/utils"
	"github.com/adshao/go-binance/v2/futures"
)

type PreparedOrder struct {
	Symbol     string
	Quantity   string
	StopPrice  string
	Side       futures.SideType
	Expiration time.Duration
}

func (bc *BinanceClient) PrepareOrder(ctx context.Context, t *models.Trade) (*PreparedOrder, error) {
	var po PreparedOrder
	q, err := bc.getQuantity(t)
	if err != nil {
		return nil, err
	}
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	stopPrice, err := bc.calculateStopPrice(t, k, s)
	if err != nil {
		return nil, err
	}
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeBuy
	} else {
		side = futures.SideTypeSell
	}
	candleDuration, err := types.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}

	candleCloseTime := utils.ConvertTime(k.CloseTime)
	remainingTime := candleDuration + time.Until(candleCloseTime)
	if remainingTime < 0 {
		return nil, fmt.Errorf("remaining time should not be negative number: %d", remainingTime)
	}

	po.Symbol = t.Symbol
	po.Quantity = q
	po.Side = side
	po.StopPrice = stopPrice
	po.Expiration = remainingTime

	return &po, nil
}

func (bc *BinanceClient) PlacePreparedOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		StopPrice(p.StopPrice).
		Type(futures.OrderTypeStopMarket).
		WorkingType(futures.WorkingTypeMarkPrice)

	return order.Do(ctx)
}

func (bc *BinanceClient) CancelOrder(ctx context.Context, orderID int64, symbol string) (*futures.CancelOrderResponse, error) {
	order := bc.client.NewCancelOrderService().
		OrderID(orderID).
		Symbol(symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) GetOrder(ctx context.Context, orderID int64, symbol string) (*futures.Order, error) {
	order := bc.client.NewGetOrderService().
		OrderID(orderID).
		Symbol(symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) PrepareStopLossOrder(ctx context.Context, t *models.Trade, ou *futures.WsOrderTradeUpdate) (*PreparedOrder, error) {
	var po PreparedOrder
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	if err != nil {
		return nil, err
	}
	// TODO check if activation Price is true, or need to change it to other price like average price
	stopPrice, err := bc.calculateStopLossPrice(t, k, s, ou.ActivationPrice)
	if err != nil {
		return nil, err
	}
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}

	po.Symbol = t.Symbol
	po.Quantity = ou.OriginalQty
	po.Side = side
	po.StopPrice = stopPrice
	po.Expiration = 0

	return &po, nil
}

func (bc *BinanceClient) PrepareTakeProfitOrder(ctx context.Context, t *models.Trade, ou *futures.WsOrderTradeUpdate) (*PreparedOrder, error) {
	var po PreparedOrder
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	if err != nil {
		return nil, err
	}
	// TODO check if activation Price is true, or need to change it to other price like average price
	stopPrice, err := bc.calculateTakeProfitPrice(t, k, s, ou.ActivationPrice)
	if err != nil {
		return nil, err
	}
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}

	po.Symbol = t.Symbol
	po.Quantity = ou.OriginalQty
	po.Side = side
	po.StopPrice = stopPrice
	po.Expiration = 0

	return &po, nil
}

func (bc *BinanceClient) PlacePreparedStopLossOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		WorkingType(futures.WorkingTypeMarkPrice).
		Type(futures.OrderTypeStopMarket).
		StopPrice(p.StopPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) PlacePreparedTakeProfitOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		WorkingType(futures.WorkingTypeMarkPrice).
		Type(futures.OrderTypeTakeProfitMarket).
		StopPrice(p.StopPrice)
	return order.Do(ctx)
}
