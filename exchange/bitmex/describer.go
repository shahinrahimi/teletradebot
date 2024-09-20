package bitmex

import (
	"context"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (mc *BitmexClient) FetchDescriber(ctx context.Context, t *models.Trade) (*models.Describer, error) {

	k, err := mc.GetLastClosedCandleOld(ctx, t)
	if err != nil {
		return nil, err
	}
	i, err := mc.GetInstrument(ctx, t)
	if err != nil {
		return nil, err
	}
	// TODO after changing value from last to mark this will be edited
	k.High = k.High * 1.005
	k.Low = k.Low * 0.995
	sp, err := t.CalculateStopPrice(k.High, k.Low)
	if err != nil {
		return nil, err
	}
	sl, err := t.CalculateStopLossPrice(k.High, k.Low, sp)
	if err != nil {
		return nil, err
	}
	tp, err := t.CalculateTakeProfitPrice(k.High, k.Low, sp)
	if err != nil {
		return nil, err
	}

	dur, err := types.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}

	candleDuration, err := types.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}

	return &models.Describer{
		OpenTime:        k.Timestamp.Add(-dur),
		CloseTime:       k.Timestamp,
		Open:            k.Open,
		Close:           k.Close,
		High:            k.High,
		Low:             k.Low,
		StopPrice:       sp,
		TakeProfitPrice: tp,
		CandleDuration:  candleDuration,
		StopLossPrice:   sl,
		TickSize:        i.TickSize,
		LotSize:         float64(i.LotSize),
	}, nil
}
