package bitmex

import (
	"context"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (mc *BitmexClient) FetchDescriber(ctx context.Context, t *models.Trade) (*models.Describer, error) {
	candleDuration, err := types.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}
	k, err := mc.GetLastClosedCandle(t.Symbol, candleDuration)
	if err != nil {
		return nil, err
	}

	// k, err := mc.GetLastClosedCandleOld(ctx, t)
	// if err != nil {
	// 	return nil, err
	// }
	i, err := mc.GetInstrument(ctx, t.Symbol)
	if err != nil {
		return nil, err
	}
	// TODO after changing value from last to mark this will be edited
	// k.High = k.High * 1.005
	// k.Low = k.Low * 0.995
	sp, err := t.CalculateStopPrice(k.High, k.Low)
	if err != nil {
		return nil, err
	}
	sl, err := t.CalculateStopLossPrice(k.High, k.Low, sp, false)
	if err != nil {
		return nil, err
	}
	tp, err := t.CalculateTakeProfitPrice(k.High, k.Low, sp, false)
	if err != nil {
		return nil, err
	}

	rsl, err := t.CalculateStopLossPrice(k.High, k.Low, sp, true)
	if err != nil {
		return nil, err
	}
	rtp, err := t.CalculateTakeProfitPrice(k.High, k.Low, sp, true)
	if err != nil {
		return nil, err
	}

	// dur, err := types.GetDuration(t.Timeframe)
	// if err != nil {
	// 	return nil, err
	// }

	return &models.Describer{
		TradeID:        t.ID,
		Symbol:         t.Symbol,
		Size:           t.Size,
		TakeProfitSize: t.TakeProfitSize,
		StopLossSize:   t.StopLossSize,
		// OpenTime:        k.Timestamp.Add(-dur),
		// CloseTime:       k.Timestamp,
		OpenTime:               k.OpenTime,
		CloseTime:              k.CloseTime,
		Open:                   k.Open,
		Close:                  k.Close,
		High:                   k.High,
		Low:                    k.Low,
		StopPrice:              sp,
		TakeProfitPrice:        tp,
		StopLossPrice:          sl,
		ReverseStopLossPrice:   rsl,
		ReverseTakeProfitPrice: rtp,
		CandleDuration:         candleDuration,
		PricePrecision:         i.TickSize,
		QuantityPrecision:      float64(i.LotSize),
	}, nil
}
