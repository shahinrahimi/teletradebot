package binance

import (
	"context"
	"fmt"
	"math"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

type PreparedOrder struct {
	Symbol    string
	Quantity  string
	StopPrice string
	Side      futures.SideType
}

func (bc *BinanceClient) prepareDescriberForMainOrder(ctx context.Context, d *models.Describer, t *models.Trade) (*PreparedOrder, error) {
	var po PreparedOrder
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeBuy
	} else {
		side = futures.SideTypeSell
	}

	balance, err := bc.getAvailableBalance()
	if err != nil {
		return nil, err
	}
	price, err := bc.getLatestPrice(t)
	if err != nil {
		return nil, err
	}

	size := balance * float64(t.Size) / 100
	quantity := size / price

	// adjust quantity based on symbol quantity precision
	quantityPrecision := math.Pow10(int(-d.QuantityPrecision))
	quantity = math.Floor(quantity/quantityPrecision) * quantityPrecision
	q := fmt.Sprintf("%.*f", d.QuantityPrecision, quantity)

	// adjust price based on symbol price precision
	pricePrecision := math.Pow10(int(-d.PricePrecision))
	stopPrice := math.Floor(d.StopPrice/pricePrecision) * pricePrecision
	p := fmt.Sprintf("%.*f", d.PricePrecision, stopPrice)

	po.Symbol = t.Symbol
	po.Side = side
	po.Quantity = q
	po.StopPrice = p

	return &po, nil
}

func (bc *BinanceClient) prepareDescriberForStopLossOrder(ctx context.Context, d *models.Describer, t *models.Trade, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
	var po PreparedOrder
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}

	// adjust price based on symbol price precision
	pricePrecision := math.Pow10(int(-d.PricePrecision))
	stopPrice := math.Floor(d.StopLossPrice/pricePrecision) * pricePrecision
	p := fmt.Sprintf("%.*f", d.PricePrecision, stopPrice)

	po.Symbol = t.Symbol
	po.Side = side
	po.Quantity = ou.OriginalQty
	po.StopPrice = p

	return &po
}

func (bc *BinanceClient) prepareDescriberForTakeProfitOrder(ctx context.Context, d *models.Describer, t *models.Trade, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
	var po PreparedOrder
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}

	// adjust price based on symbol price precision
	pricePrecision := math.Pow10(int(-d.PricePrecision))
	stopPrice := math.Floor(d.TakeProfitPrice/pricePrecision) * pricePrecision
	p := fmt.Sprintf("%.*f", d.PricePrecision, stopPrice)

	po.Symbol = t.Symbol
	po.Side = side
	po.Quantity = ou.OriginalQty
	po.StopPrice = p

	return &po
}

// func (bc *BinanceClient) prepareOrder(ctx context.Context, t *models.Trade) (*PreparedOrder, error) {
// 	var po PreparedOrder
// 	q, err := bc.getQuantity(t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	k, err := bc.getLastClosedKline(ctx, t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	s, err := bc.getSymbol(t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	stopPrice, err := bc.calculateStopPrice(t, k, s)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var side futures.SideType
// 	if t.Side == types.SIDE_L {
// 		side = futures.SideTypeBuy
// 	} else {
// 		side = futures.SideTypeSell
// 	}
// 	candleDuration, err := types.GetDuration(t.Timeframe)
// 	if err != nil {
// 		return nil, err
// 	}

// 	candleCloseTime := utils.ConvertTime(k.CloseTime)
// 	remainingTime := candleDuration + time.Until(candleCloseTime)
// 	if remainingTime < 0 {
// 		return nil, fmt.Errorf("remaining time should not be negative number: %d", remainingTime)
// 	}

// 	po.Symbol = t.Symbol
// 	po.Quantity = q
// 	po.Side = side
// 	po.StopPrice = stopPrice
// 	po.Expiration = remainingTime

// 	return &po, nil
// }

// func (bc *BinanceClient) prepareSLOrder(ctx context.Context, t *models.Trade, ou *futures.WsOrderTradeUpdate) (*PreparedOrder, error) {
// 	var po PreparedOrder
// 	k, err := bc.getLastClosedKline(ctx, t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	s, err := bc.getSymbol(t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// TODO check if activation Price is true, or need to change it to other price like average price
// 	stopPrice, err := bc.calculateStopLossPrice(t, k, s, ou.StopPrice)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var side futures.SideType
// 	if t.Side == types.SIDE_L {
// 		side = futures.SideTypeSell
// 	} else {
// 		side = futures.SideTypeBuy
// 	}

// 	po.Symbol = t.Symbol
// 	po.Quantity = ou.OriginalQty
// 	po.Side = side
// 	po.StopPrice = stopPrice
// 	po.Expiration = 0

// 	return &po, nil
// }

// func (bc *BinanceClient) prepareTPOrder(ctx context.Context, t *models.Trade, ou *futures.WsOrderTradeUpdate) (*PreparedOrder, error) {
// 	var po PreparedOrder
// 	k, err := bc.getLastClosedKline(ctx, t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	s, err := bc.getSymbol(t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// TODO check if activation Price is true, or need to change it to other price like average price
// 	stopPrice, err := bc.calculateTakeProfitPrice(t, k, s, ou.StopPrice)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var side futures.SideType
// 	if t.Side == types.SIDE_L {
// 		side = futures.SideTypeSell
// 	} else {
// 		side = futures.SideTypeBuy
// 	}

// 	po.Symbol = t.Symbol
// 	po.Quantity = ou.OriginalQty
// 	po.Side = side
// 	po.StopPrice = stopPrice
// 	po.Expiration = 0

// 	return &po, nil
// }
