package binance

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

type PreparedOrder struct {
	Symbol     string
	Quantity   string
	StopPrice  string
	Side       futures.SideType
	Expiration time.Duration
}

func (bc *BinanceClient) prepareOrder(ctx context.Context, t *models.Trade) (*PreparedOrder, error) {
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
	if err != nil {
		return nil, err
	}
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

func (bc *BinanceClient) prepareSLOrder(ctx context.Context, t *models.Trade, ou *futures.WsOrderTradeUpdate) (*PreparedOrder, error) {
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
	stopPrice, err := bc.calculateStopLossPrice(t, k, s, ou.StopPrice)
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

func (bc *BinanceClient) prepareTPOrder(ctx context.Context, t *models.Trade, ou *futures.WsOrderTradeUpdate) (*PreparedOrder, error) {
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
	stopPrice, err := bc.calculateTakeProfitPrice(t, k, s, ou.StopPrice)
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
