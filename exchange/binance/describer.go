package binance

import (
	"context"
	"strconv"
	"time"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) FetchDescriber(ctx context.Context, t *models.Trade) (*models.Describer, error) {
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	high, err := strconv.ParseFloat(k.High, 64)
	if err != nil {
		return nil, err
	}
	low, err := strconv.ParseFloat(k.Low, 64)
	if err != nil {
		return nil, err
	}
	open, err := strconv.ParseFloat(k.Open, 64)
	if err != nil {
		return nil, err
	}
	close, err := strconv.ParseFloat(k.Close, 64)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	if err != nil {
		return nil, err
	}
	sp, err := t.CalculateStopPrice(high, low)
	if err != nil {
		return nil, err
	}
	sl, err := t.CalculateStopLossPrice(high, low, sp, false)
	if err != nil {
		return nil, err
	}
	tp, err := t.CalculateTakeProfitPrice(high, low, sp, false)
	if err != nil {
		return nil, err
	}
	rsl, err := t.CalculateStopLossPrice(high, low, sp, true)
	if err != nil {
		return nil, err
	}
	rtp, err := t.CalculateTakeProfitPrice(high, low, sp, true)
	if err != nil {
		return nil, err
	}

	candleDuration, err := types.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}

	return &models.Describer{
		OpenTime:               utils.ConvertTime(k.OpenTime),
		CloseTime:              utils.ConvertTime(k.CloseTime).Add(time.Second),
		Open:                   open,
		Close:                  close,
		High:                   high,
		Low:                    low,
		StopPrice:              sp,
		StopLossPrice:          sl,
		TakeProfitPrice:        tp,
		ReverseStopLossPrice:   rsl,
		ReverseTakeProfitPrice: rtp,
		CandleDuration:         candleDuration,
		PricePrecision:         s.PricePrecision,
		QuantityPrecision:      s.QuantityPrecision,
	}, nil
}
