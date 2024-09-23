package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) fetchBalance(ctx context.Context) (float64, error) {
	res, err := bc.Client.NewGetBalanceService().Do(ctx)
	if err != nil {
		return 0, err
	}
	for _, balance := range res {
		if balance.Asset == "USDT" {
			return strconv.ParseFloat(balance.Balance, 64)
		}
	}
	return 0, fmt.Errorf("Asset USDT not found")
}

func (bc *BinanceClient) fetchPrice(ctx context.Context, symbol string) (float64, error) {
	res, err := bc.Client.NewListPricesService().Do(ctx)
	if err != nil {
		return 0, err
	}
	for _, sp := range res {
		if sp.Symbol == symbol {
			return strconv.ParseFloat(sp.Price, 64)
		}
	}
	return 0, fmt.Errorf("Symbol %s not found", symbol)
}

func (bc *BinanceClient) fetchLastCompletedCandle(ctx context.Context, symbol string, t models.TimeframeType) (*futures.Kline, error) {
	klines, err := bc.Client.NewMarkPriceKlinesService().
		Limit(100).
		Interval(string(t)).
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	// loop through klines and return the most recent completely closed candle
	for i := len(klines) - 1; i >= 0; i-- {
		candleCloseTime := utils.ConvertTime(klines[i].CloseTime)
		// check if close time in the past
		if (time.Until(candleCloseTime)) < 0 {
			return klines[i], nil
		}
	}

	return nil, fmt.Errorf("failed to locate before last candle")
}

func (bc *BinanceClient) FetchDescriber(ctx context.Context, t *models.Trade) (*models.Describer, error) {
	timeframeDur, err := models.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}
	k, err := bc.fetchLastCompletedCandle(ctx, t.Symbol, t.Timeframe)
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
	s, err := bc.GetSymbol(t.Symbol)
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

	return &models.Describer{
		TradeID:           t.ID,
		Symbol:            t.Symbol,
		Size:              t.Size,
		TakeProfitSize:    t.TakeProfitSize,
		StopLossSize:      t.StopLossSize,
		ReverseMultiplier: t.ReverseMultiplier,
		TimeFrame:         timeframeDur,

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
		PricePrecision:         float64(s.PricePrecision),
		QuantityPrecision:      float64(s.QuantityPrecision),
	}, nil
}
