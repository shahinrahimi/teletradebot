package bitmex

import (
	"context"
	"time"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

type PreparedOrder struct {
	Symbol     string
	Quantity   float32
	StopPrice  float64
	Side       string
	Expiration time.Duration
}

func (mc *BitmexClient) prepareOrder(ctx context.Context, t *models.Trade) (*PreparedOrder, error) {
	var p PreparedOrder

	// timeframeDur, err := types.GetDuration(t.Timeframe)
	// if err != nil {
	// 	return nil, err
	// }
	candle, err := mc.GetLastClosedCandleOld(ctx, t)
	if err != nil {
		mc.l.Printf("Error fetching last closed candle: %v", err)
		return nil, err
	}

	if t.Side == types.SIDE_L {
		p.Side = SideTypeBuy
	} else {
		p.Side = SideTypeSell
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
	stopPrice, err := mc.calculateStopPrice(t, candle.High, candle.Low, instrument.TickSize)
	if err != nil {
		return nil, err
	}

	quantity, err := mc.calculateQuantity(t, balance, instrument.MarkPrice, float64(instrument.LotSize))
	if err != nil {
		return nil, err
	}

	expiration, err := mc.calculateExpiration(t, candle.Timestamp)
	if err != nil {
		return nil, err
	}

	p.Quantity = float32(quantity)
	p.StopPrice = stopPrice
	p.Expiration = expiration
	p.Symbol = t.Symbol

	return &p, nil
}

func (mc *BitmexClient) prepareSLOrder(ctx context.Context, t *models.Trade) (*PreparedOrder, error) {
	return nil, nil

	// timeframeDur, err := types.GetDuration(t.Timeframe)
	// if err != nil {
	// 	return nil, err
	// }
	// candle, err := mc.GetLastClosedCandleOld(ctx, t)
	// if err != nil {
	// 	mc.l.Printf("Error fetching last closed candle: %v", err)
	// 	return nil, err
	// }

	// if t.Side == types.SIDE_L {
	// 	p.Side = SideTypeSell
	// } else {
	// 	p.Side = SideTypeBuy
	// }

	// instrument, err := mc.GetInstrument(ctx, t)
	// if err != nil {
	// 	mc.l.Printf("Error fetching instrument: %v", err)
	// 	return nil, err
	// }
	// stopPrice, err := mc.calculateStopLossPrice(t, candle.High, candle.Low, instrument.TickSize)
	// if err != nil {
	// 	return nil, err
	// }

	// quantity, err := mc.calculateQuantity(t, balance, instrument.MarkPrice, float64(instrument.LotSize))
	// if err != nil {
	// 	return nil, err
	// }

	// expiration, err := mc.calculateExpiration(t, candle.Timestamp)
	// if err != nil {
	// 	return nil, err
	// }

	// p.Quantity = float32(quantity)
	// p.StopPrice = stopPrice
	// p.Expiration = expiration
	// p.Symbol = t.Symbol

	// return &p, nil
}

func (nc *BitmexClient) prepareTPOrder(ctx context.Context, t *models.Trade) (*PreparedOrder, error) {
	return nil, nil
}
