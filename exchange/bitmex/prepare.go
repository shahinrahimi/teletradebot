package bitmex

import (
	"context"
	"fmt"
	"math"

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

func (mc *BitmexClient) prepareDescriberForMainOrder(ctx context.Context, d *models.Describer, t *models.Trade) (*PreparedOrder, error) {
	var po PreparedOrder
	var side SideType
	if t.Side == types.SIDE_L {
		side = SideTypeBuy
	} else {
		side = SideTypeSell
	}

	balance, err := mc.GetBalanceUSDt(ctx)
	if err != nil {
		mc.l.Printf("Error fetching balance: %v", err)
		return nil, err
	}

	instrument, err := mc.GetInstrument(ctx, t)
	if err != nil {
		mc.l.Printf("Error fetching instrument: %v", err)
		return nil, err
	}

	contractSize, exist := config.ContractSizes[t.Symbol]
	if !exist {
		return nil, fmt.Errorf("contract size not found for symbol %s", t.Symbol)
	}

	balance = balance / 1000000 // balance in USDT
	price := instrument.MarkPrice

	size := balance * (float64(t.Size) / 100)
	quantity := size / (price * contractSize)
	// adjust quantity based on symbol lot size
	q := math.Floor(quantity/d.LotSize) * d.LotSize
	// log all variable for debug
	mc.l.Printf("balance: %f, price: %f, size: %f, quantity: %f, lotsize: %f, calculated quantity: %f, maximum quantity: %f", balance, price, size, quantity, d.LotSize, q, instrument.MaxOrderQty)

	if q < d.LotSize {
		return nil, fmt.Errorf("the calculated quantity (%.2f) less than the lotsize (%.1f)", q, d.LotSize)
	}
	if q > float64(instrument.MaxOrderQty) {
		q = float64(instrument.MaxOrderQty)
		msg := fmt.Sprintf("the calculated quantity (%.2f) greater than the maximum order quantity (%.2f). so the calculated quantity will be adjusted to %.2f", q, instrument.MaxOrderQty, instrument.MaxOrderQty)
		mc.l.Println(msg)
		//return nil, fmt.Errorf("the calculated quantity (%.2f) greater than the maximum order quantity (%.2f)", q, instrument.MaxOrderQty)
	}

	// adjust price based on symbol tick size
	ticks := d.StopPrice / d.TickSize
	p := math.Round(ticks) * d.TickSize

	// log calculated price
	mc.l.Printf("calculated price: %f for the symbol: %s", p, t.Symbol)

	po.Symbol = t.Symbol
	po.Side = side
	po.Quantity = q
	po.StopPrice = p

	return &po, nil
}

func (mc *BitmexClient) prepareDescriberForStopLossOrder(ctx context.Context, d *models.Describer, t *models.Trade, od *OrderData) *PreparedOrder {
	var po PreparedOrder
	var side SideType
	if t.Side == types.SIDE_L {
		side = SideTypeSell
	} else {
		side = SideTypeBuy
	}

	// adjust price based on symbol tick size
	ticks := d.StopLossPrice / d.TickSize
	p := math.Round(ticks) * d.TickSize

	po.Symbol = t.Symbol
	po.Side = side
	// TODO check if OrderQty is the correct value
	po.Quantity = float64(od.OrderQty)
	po.StopPrice = p

	return &po
}

func (mc *BitmexClient) prepareDescriberForTakeProfitOrder(ctx context.Context, d *models.Describer, t *models.Trade, od *OrderData) *PreparedOrder {
	var po PreparedOrder
	var side SideType
	if t.Side == types.SIDE_L {
		side = SideTypeSell
	} else {
		side = SideTypeBuy
	}

	// adjust price based on symbol tick size
	ticks := d.TakeProfitPrice / d.TickSize
	p := math.Round(ticks) * d.TickSize

	po.Symbol = t.Symbol
	po.Side = side
	// TODO check if OrderQty is the correct value
	po.Quantity = float64(od.OrderQty)
	po.StopPrice = p

	return &po
}

// func (mc *BitmexClient) prepareOrder(ctx context.Context, d *models.Describer, t *models.Trade) (*PreparedOrder, error) {
// 	var p PreparedOrder
// 	var side SideType
// 	if t.Side == types.SIDE_L {
// 		side = SideTypeBuy
// 	} else {
// 		side = SideTypeSell
// 	}

// 	balance, err := mc.GetBalanceUSDt(ctx)
// 	if err != nil {
// 		mc.l.Printf("Error fetching balance: %v", err)
// 		return nil, err
// 	}

// 	instrument, err := mc.GetInstrument(ctx, t)
// 	if err != nil {
// 		mc.l.Printf("Error fetching instrument: %v", err)
// 		return nil, err
// 	}
// 	price := instrument.MarkPrice

// 	candle, err := mc.GetLastClosedCandleOld(ctx, t)
// 	if err != nil {
// 		mc.l.Printf("Error fetching last closed candle: %v", err)
// 		return nil, err
// 	}

// 	// TODO will change after solution on the market price
// 	stopPrice, err := mc.calculateStopPrice(t, candle.High*1.005, candle.Low*0.995, instrument.TickSize)
// 	if err != nil {
// 		return nil, err
// 	}

// 	quantity, err := mc.calculateQuantity(t, balance, instrument.MarkPrice, float64(instrument.LotSize))
// 	if err != nil {
// 		return nil, err
// 	}

// 	expiration, err := mc.calculateExpiration(t, candle.Timestamp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if t.Side == types.SIDE_L {
// 		p.Side = SideTypeBuy
// 	} else {
// 		p.Side = SideTypeSell
// 	}

// 	p.StopPrice = stopPrice
// 	p.Quantity = float32(quantity)
// 	p.Expiration = expiration
// 	p.Symbol = t.Symbol

// 	return &p, nil

// }
