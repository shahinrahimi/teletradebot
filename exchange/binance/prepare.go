package binance

import (
	"context"
	"strconv"

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
	balanceChan := make(chan float64)
	priceChan := make(chan float64)
	errChan := make(chan error, 2)

	go func() {
		balance, err := bc.fetchBalance(ctx)
		if err != nil {
			errChan <- err
			return
		}
		balanceChan <- balance
	}()

	go func() {
		price, err := bc.fetchPrice(ctx, d.Symbol)
		if err != nil {
			errChan <- err
			return
		}
		priceChan <- price
	}()

	var balance float64
	var price float64
	var err error

	select {
	case balance = <-balanceChan:
	case err = <-errChan:
		return nil, err
	}
	select {
	case price = <-priceChan:
	case err = <-errChan:
		return nil, err
	}

	var po PreparedOrder
	var side futures.SideType
	if d.Side == types.SIDE_L {
		side = futures.SideTypeBuy
	} else {
		side = futures.SideTypeSell
	}

	size := balance * float64(d.Size) / 100
	quantity := size / price

	p := d.AdjustPriceForBinance(d.StopPrice)
	q := d.AdjustQuantityForBinance(quantity)

	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = q
	po.StopPrice = p

	return &po, nil
}

func (bc *BinanceClient) prepareStopLossOrder(d *models.Describer, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
	var po PreparedOrder
	var side futures.SideType
	if d.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}
	p := d.AdjustPriceForBinance(d.StopLossPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = ou.OriginalQty
	po.StopPrice = p

	return &po
}

func (bc *BinanceClient) prepareTakeProfitOrder(d *models.Describer, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
	var po PreparedOrder
	var side futures.SideType
	if d.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}
	p := d.AdjustPriceForBinance(d.TakeProfitPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = ou.OriginalQty
	po.StopPrice = p
	return &po
}

func (bc *BinanceClient) prepareReverseMainOrder(d *models.Describer, ou *futures.WsOrderTradeUpdate) (*PreparedOrder, error) {
	var po PreparedOrder
	var side futures.SideType
	if d.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}
	p := d.AdjustPriceForBinance(d.StopLossPrice)
	originalQty, err := strconv.ParseFloat(ou.OriginalQty, 64)
	if err != nil {
		return nil, err
	}
	quantity := originalQty * float64(d.ReverseMultiplier)
	q := d.AdjustQuantityForBinance(quantity)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = q
	po.StopPrice = p
	return &po, err
}

func (bc *BinanceClient) prepareReverseStopLossOrder(d *models.Describer, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
	var po PreparedOrder
	var side futures.SideType
	if d.Side == types.SIDE_L {
		side = futures.SideTypeBuy
	} else {
		side = futures.SideTypeSell
	}
	p := d.AdjustPriceForBinance(d.ReverseStopLossPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = ou.OriginalQty
	po.StopPrice = p
	return &po
}

func (bc *BinanceClient) prepareReverseTakeProfitOrder(d *models.Describer, ou *futures.WsOrderTradeUpdate) *PreparedOrder {
	var po PreparedOrder
	var side futures.SideType
	if d.Side == types.SIDE_L {
		side = futures.SideTypeBuy
	} else {
		side = futures.SideTypeSell
	}
	p := d.AdjustPriceForBinance(d.ReverseTakeProfitPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = ou.OriginalQty
	po.StopPrice = p
	return &po
}
