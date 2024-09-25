package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) fetchBalance(ctx context.Context) (float64, error) {
	res, err := bc.client.NewGetBalanceService().Do(ctx)
	if err != nil {
		return 0, err
	}
	for _, balance := range res {
		if balance.Asset == "USDT" {
			return strconv.ParseFloat(balance.Balance, 64)
		}
	}
	return 0, fmt.Errorf("asset USDT not found")
}

func (bc *BinanceClient) fetchPrice(ctx context.Context, symbol string) (float64, error) {
	res, err := bc.client.NewListPricesService().Do(ctx)
	if err != nil {
		return 0, err
	}
	for _, sp := range res {
		if sp.Symbol == symbol {
			return strconv.ParseFloat(sp.Price, 64)
		}
	}
	return 0, fmt.Errorf("symbol %s not found", symbol)
}

func (bc *BinanceClient) fetchLastCompletedCandle(ctx context.Context, symbol string, t models.TimeframeType) (*futures.Kline, error) {
	klines, err := bc.client.NewMarkPriceKlinesService().
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

func (bc *BinanceClient) FetchInterpreter(ctx context.Context, t *models.Trade) (*models.Interpreter, error) {
	errChan := make(chan error, 2)
	balanceChan := make(chan float64, 1)
	priceChan := make(chan float64, 1)
	candleChan := make(chan *futures.Kline, 1)
	symbolChan := make(chan *futures.Symbol, 1)

	go func() {
		candle, err := bc.fetchLastCompletedCandle(ctx, t.Symbol, t.Timeframe)
		if err != nil {
			errChan <- err
		}
		candleChan <- candle
	}()
	go func() {
		symbol, err := bc.GetSymbol(t.Symbol)
		if err != nil {
			errChan <- err
		}
		symbolChan <- symbol
	}()
	go func() {
		balance, err := bc.fetchBalance(ctx)
		if err != nil {
			errChan <- err
		}
		balanceChan <- balance
	}()
	go func() {
		price, err := bc.fetchPrice(ctx, t.Symbol)
		if err != nil {
			errChan <- err
		}
		priceChan <- price
	}()

	var candle *futures.Kline
	var symbol *futures.Symbol
	var balance float64
	var price float64
	var err error

	for i := 0; i < 4; i++ {
		select {
		case candle = <-candleChan:
		case err = <-errChan:
		case balance = <-balanceChan:
		case price = <-priceChan:
		case symbol = <-symbolChan:
		}
		if err != nil {
			return nil, err
		}
	}

	timeframeDur, err := models.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}
	high, err := strconv.ParseFloat(candle.High, 64)
	if err != nil {
		return nil, err
	}
	low, err := strconv.ParseFloat(candle.Low, 64)
	if err != nil {
		return nil, err
	}
	open, err := strconv.ParseFloat(candle.Open, 64)
	if err != nil {
		return nil, err
	}
	close, err := strconv.ParseFloat(candle.Close, 64)
	if err != nil {
		return nil, err
	}
	ep, err := t.CalculateEntryPrice(high, low)
	if err != nil {
		return nil, err
	}
	sl, err := t.CalculateStopLossPrice(high, low, ep, false)
	if err != nil {
		return nil, err
	}
	tp, err := t.CalculateTakeProfitPrice(high, low, ep, false)
	if err != nil {
		return nil, err
	}
	rsl, err := t.CalculateStopLossPrice(high, low, sl, true)
	if err != nil {
		return nil, err
	}
	rtp, err := t.CalculateTakeProfitPrice(high, low, sl, true)
	if err != nil {
		return nil, err
	}

	size := balance * float64(t.Size) / 100
	quantity := size / price
	rQuantity := quantity * float64(t.ReverseMultiplier)

	return &models.Interpreter{
		Balance:         balance,
		Price:           price,
		Quantity:        quantity,
		ReverseQuantity: rQuantity,
		Exchange:        types.ExchangeBinance,

		TradeID:           t.ID,
		Symbol:            t.Symbol,
		Size:              t.Size,
		Side:              t.Side,
		TakeProfitSize:    t.TakeProfitSize,
		StopLossSize:      t.StopLossSize,
		ReverseMultiplier: t.ReverseMultiplier,
		TimeFrameDur:      timeframeDur,

		OpenTime:               utils.ConvertTime(candle.OpenTime),
		CloseTime:              utils.ConvertTime(candle.CloseTime).Add(time.Second),
		Open:                   open,
		Close:                  close,
		High:                   high,
		Low:                    low,
		EntryPrice:             ep,
		StopLossPrice:          sl,
		TakeProfitPrice:        tp,
		ReverseEntryPrice:      sl,
		ReverseStopLossPrice:   rsl,
		ReverseTakeProfitPrice: rtp,
		PricePrecision:         symbol.PricePrecision,
		QuantityPrecision:      symbol.QuantityPrecision,
	}, nil
}