package binance

import (
	"context"

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

func (bc *BinanceClient) prepareMainOrder(ctx context.Context, d *models.Describer) (*PreparedOrder, error) {
	var po PreparedOrder
	var side futures.SideType
	if d.Side == types.SIDE_L {
		side = futures.SideTypeBuy
	} else {
		side = futures.SideTypeSell
	}

	balance, err := bc.getAvailableBalance()
	if err != nil {
		return nil, err
	}
	price, err := bc.getLatestPrice(d.Symbol)
	if err != nil {
		return nil, err
	}

	size := balance * float64(d.Size) / 100
	quantity := size / price

	p := d.GetValueWithPricePrecisionString(d.StopPrice)
	q := d.GetValueWithQuantityPrecisionString(quantity)

	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = q
	po.StopPrice = p

	return &po, nil
}

func (bc *BinanceClient) prepareStopLossOrder(ctx context.Context, d *models.Describer, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
	var po PreparedOrder
	var side futures.SideType
	if d.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}
	p := d.GetValueWithPricePrecisionString(d.StopLossPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = ou.OriginalQty
	po.StopPrice = p

	return &po
}

func (bc *BinanceClient) prepareTakeProfitOrder(ctx context.Context, d *models.Describer, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
	var po PreparedOrder
	var side futures.SideType
	if d.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}
	p := d.GetValueWithPricePrecisionString(d.TakeProfitPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = ou.OriginalQty
	po.StopPrice = p
	return &po
}

// func (bc *BinanceClient) prepareDescriberForMainOrder(ctx context.Context, d *models.Describer, t *models.Trade) (*PreparedOrder, error) {
// 	var po PreparedOrder
// 	var side futures.SideType
// 	if t.Side == types.SIDE_L {
// 		side = futures.SideTypeBuy
// 	} else {
// 		side = futures.SideTypeSell
// 	}

// 	balance, err := bc.getAvailableBalance()
// 	if err != nil {
// 		return nil, err
// 	}
// 	price, err := bc.getLatestPrice(t)
// 	if err != nil {
// 		return nil, err
// 	}

// 	size := balance * float64(t.Size) / 100
// 	quantity := size / price

// 	// adjust quantity based on symbol quantity precision
// 	quantityPrecision := math.Pow10(int(-d.QuantityPrecision))
// 	quantity = math.Floor(quantity/quantityPrecision) * quantityPrecision
// 	q := fmt.Sprintf("%.*f", d.QuantityPrecision, quantity)

// 	// adjust price based on symbol price precision
// 	pricePrecision := math.Pow10(int(-d.PricePrecision))
// 	stopPrice := math.Floor(d.StopPrice/pricePrecision) * pricePrecision
// 	p := fmt.Sprintf("%.*f", d.PricePrecision, stopPrice)

// 	po.Symbol = t.Symbol
// 	po.Side = side
// 	po.Quantity = q
// 	po.StopPrice = p

// 	return &po, nil
// }

// func (bc *BinanceClient) prepareDescriberForStopLossOrder(ctx context.Context, d *models.Describer, t *models.Trade, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
// 	var po PreparedOrder
// 	var side futures.SideType
// 	if t.Side == types.SIDE_L {
// 		side = futures.SideTypeSell
// 	} else {
// 		side = futures.SideTypeBuy
// 	}

// 	// adjust price based on symbol price precision
// 	pricePrecision := math.Pow10(int(-d.PricePrecision))
// 	stopPrice := math.Floor(d.StopLossPrice/pricePrecision) * pricePrecision
// 	p := fmt.Sprintf("%.*f", d.PricePrecision, stopPrice)

// 	po.Symbol = t.Symbol
// 	po.Side = side
// 	po.Quantity = ou.OriginalQty
// 	po.StopPrice = p

// 	return &po
// }

// func (bc *BinanceClient) prepareDescriberForTakeProfitOrder(ctx context.Context, d *models.Describer, t *models.Trade, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
// 	var po PreparedOrder
// 	var side futures.SideType
// 	if t.Side == types.SIDE_L {
// 		side = futures.SideTypeSell
// 	} else {
// 		side = futures.SideTypeBuy
// 	}

// 	// adjust price based on symbol price precision
// 	pricePrecision := math.Pow10(int(-d.PricePrecision))
// 	stopPrice := math.Floor(d.TakeProfitPrice/pricePrecision) * pricePrecision
// 	p := fmt.Sprintf("%.*f", d.PricePrecision, stopPrice)

// 	po.Symbol = t.Symbol
// 	po.Side = side
// 	po.Quantity = ou.OriginalQty
// 	po.StopPrice = p

// 	return &po
// }

// func (bc *BinanceClient) prepareDescriberForReverseStopLossOrder(ctx context.Context, d *models.Describer, t *models.Trade) *PreparedOrder {
// 	var po PreparedOrder
// 	var side futures.SideType
// 	if t.Side == types.SIDE_L {
// 		side = futures.SideTypeBuy
// 	} else {
// 		side = futures.SideTypeSell
// 	}

// 	// adjust price based on symbol price precision
// 	pricePrecision := math.Pow10(int(-d.PricePrecision))
// 	stopPrice := math.Floor(d.ReverseStopLossPrice/pricePrecision) * pricePrecision
// 	p := fmt.Sprintf("%.*f", d.PricePrecision, stopPrice)

// 	po.Symbol = t.Symbol
// 	po.Side = side
// 	po.Quantity = "0" //ou.OriginalQty
// 	po.StopPrice = p

// 	return &po
// }

// func (bc *BinanceClient) prepareDescriberForReverseTakeProfitOrder(ctx context.Context, d *models.Describer, t *models.Trade, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
// 	var po PreparedOrder
// 	var side futures.SideType
// 	if t.Side == types.SIDE_L {
// 		side = futures.SideTypeBuy
// 	} else {
// 		side = futures.SideTypeSell
// 	}
// 	// adjust price based on symbol price precision
// 	pricePrecision := math.Pow10(int(-d.PricePrecision))
// 	stopPrice := math.Floor(d.ReverseTakeProfitPrice/pricePrecision) * pricePrecision
// 	p := fmt.Sprintf("%.*f", d.PricePrecision, stopPrice)

// 	po.Symbol = t.Symbol
// 	po.Side = side
// 	po.Quantity = ou.OriginalQty
// 	po.StopPrice = p

// 	return &po
// }
