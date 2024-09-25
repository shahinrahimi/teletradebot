package bitmexold

import (
	"context"
	"fmt"

	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

type PreparedOrder struct {
	Symbol    string
	Quantity  float64
	StopPrice float64
	Side      SideType
}

func (mc *BitmexClient) prepareMainOrder(ctx context.Context, d *models.Describer) (*PreparedOrder, error) {

	balanceChan := make(chan float64)
	priceChan := make(chan float64)
	errChan := make(chan error, 2)

	go func() {
		balance, err := mc.fetchBalance(ctx)
		if err != nil {
			errChan <- err
			return
		}
		balanceChan <- balance
	}()

	go func() {
		price, err := mc.fetchPrice(ctx, d.Symbol)
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
	var side SideType
	if d.Side == types.SIDE_L {
		side = SideTypeBuy
	} else {
		side = SideTypeSell
	}

	contractSize, exist := config.ContractSizes[d.Symbol]
	if !exist {
		return nil, fmt.Errorf("contract size not found for symbol %s", d.Symbol)
	}

	balance = balance / 1000000 // balance in USDT

	size := balance * (float64(d.Size) / 100)
	quantity := size / (price * contractSize)

	p := d.AdjustPriceForBitmex(d.StopPrice)
	q := d.AdjustQuantityForBitmex(quantity)

	if q < d.LotSize {
		return nil, fmt.Errorf("the calculated quantity (%.2f) less than the lotsize (%.1f)", q, d.LotSize)
	}
	if q > d.MaxOrderQty {
		return nil, fmt.Errorf("the calculated quantity (%.2f) greater than the max order quantity (%.1f)", q, d.MaxOrderQty)
	}

	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = q
	po.StopPrice = p

	return &po, nil
}

func (mc *BitmexClient) prepareStopLossOrder(d *models.Describer, od *OrderData) *PreparedOrder {
	var po PreparedOrder
	var side SideType
	if d.Side == types.SIDE_L {
		side = SideTypeSell
	} else {
		side = SideTypeBuy
	}
	p := d.AdjustPriceForBitmex(d.StopLossPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = float64(od.OrderQty)
	po.StopPrice = p
	return &po
}

func (mc *BitmexClient) prepareTakeProfitOrder(d *models.Describer, od *OrderData) *PreparedOrder {
	var po PreparedOrder
	var side SideType
	if d.Side == types.SIDE_L {
		side = SideTypeSell
	} else {
		side = SideTypeBuy
	}
	p := d.AdjustPriceForBitmex(d.TakeProfitPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = float64(od.OrderQty)
	po.StopPrice = p
	return &po
}

func (mc *BitmexClient) prepareReverseMainOrder(d *models.Describer, od *OrderData) *PreparedOrder {
	var po PreparedOrder
	var side SideType
	if d.Side == types.SIDE_L {
		side = SideTypeSell
	} else {
		side = SideTypeBuy
	}
	p := d.AdjustPriceForBitmex(d.StopLossPrice)
	quantity := float64(od.OrderQty) * float64(d.ReverseMultiplier)
	q := d.AdjustQuantityForBitmex(quantity)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = q
	po.StopPrice = p
	return &po
}

func (mc *BitmexClient) prepareReverseStopLossOrder(d *models.Describer, od *OrderData) *PreparedOrder {
	var po PreparedOrder
	var side SideType
	if d.Side == types.SIDE_L {
		side = SideTypeBuy
	} else {
		side = SideTypeSell
	}
	p := d.AdjustPriceForBitmex(d.ReverseStopLossPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = float64(od.OrderQty)
	po.StopPrice = p
	return &po
}

func (mc *BitmexClient) prepareReverseTakeProfitOrder(d *models.Describer, od *OrderData) *PreparedOrder {
	var po PreparedOrder
	var side SideType
	if d.Side == types.SIDE_L {
		side = SideTypeBuy
	} else {
		side = SideTypeSell
	}
	p := d.AdjustPriceForBitmex(d.ReverseTakeProfitPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = float64(od.OrderQty)
	po.StopPrice = p
	return &po
}
