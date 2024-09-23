package bitmex

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
	var po PreparedOrder
	var side SideType
	if d.Side == types.SIDE_L {
		side = SideTypeBuy
	} else {
		side = SideTypeSell
	}

	balance, err := mc.fetchBalance(ctx)
	if err != nil {
		mc.l.Printf("Error fetching balance: %v", err)
		return nil, err
	}

	instrument, err := mc.fetchInstrument(ctx, d.Symbol)
	if err != nil {
		mc.l.Printf("Error fetching instrument: %v", err)
		return nil, err
	}

	price, err := mc.fetchPrice(ctx, d.Symbol)
	if err != nil {
		mc.l.Printf("Error fetching price: %v", err)
		return nil, err
	}

	contractSize, exist := config.ContractSizes[d.Symbol]
	if !exist {
		return nil, fmt.Errorf("contract size not found for symbol %s", d.Symbol)
	}

	balance = balance / 1000000 // balance in USDT

	size := balance * (float64(d.Size) / 100)
	quantity := size / (price * contractSize)

	p := d.GetValueWithPricePrecision(d.StopPrice)
	q := d.GetValueWithQuantityPrecision(quantity)

	if q < float64(instrument.LotSize) {
		return nil, fmt.Errorf("the calculated quantity (%.2f) less than the lotsize (%.1f)", q, instrument.LotSize)
	}
	if q > float64(instrument.MaxOrderQty) {
		return nil, fmt.Errorf("the calculated quantity (%.2f) greater than the max order quantity (%.1f)", q, instrument.MaxOrderQty)
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
	p := d.GetValueWithPricePrecision(d.StopLossPrice)
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
	p := d.GetValueWithPricePrecision(d.TakeProfitPrice)
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
	p := d.GetValueWithPricePrecision(d.StopLossPrice)
	quantity := float64(od.OrderQty) * float64(d.ReverseMultiplier)
	q := d.GetValueWithQuantityPrecision(quantity)
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
	p := d.GetValueWithPricePrecision(d.ReverseStopLossPrice)
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
	p := d.GetValueWithPricePrecision(d.ReverseTakeProfitPrice)
	po.Symbol = d.Symbol
	po.Side = side
	po.Quantity = float64(od.OrderQty)
	po.StopPrice = p
	return &po
}
