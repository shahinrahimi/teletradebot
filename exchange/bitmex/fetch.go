package bitmex

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"github.com/shahinrahimi/teletradebot/models"
	swagger "github.com/shahinrahimi/teletradebot/swagger"
)

func (mc *BitmexClient) fetchMargins(ctx context.Context) ([]swagger.Margin, error) {
	ctx = mc.getAuthContext(ctx)
	opts := swagger.UserApiUserGetMarginOpts{
		Currency: optional.NewString("all"),
	}
	margins, _, err := mc.client.UserApi.UserGetMargins(ctx, &opts)
	if err != nil {
		mc.l.Printf("failed to retrieve margins: %v", err)
		return nil, err
	}
	return margins, nil
}

func (mc *BitmexClient) fetchBalance(ctx context.Context) (float64, error) {
	ctx = mc.getAuthContext(ctx)
	opts := swagger.UserApiUserGetMarginOpts{
		Currency: optional.NewString("all"),
	}
	margins, _, err := mc.client.UserApi.UserGetMargins(ctx, &opts)
	if err != nil {
		mc.l.Printf("failed to retrieve margins: %v", err)
		return 0, err
	}
	for _, m := range margins {
		if m.Currency == "USDt" {
			return float64(m.AvailableMargin), nil
		}
	}
	return 0, fmt.Errorf("the currency 'USTt' not found")
}

func (mc *BitmexClient) fetchPrice(ctx context.Context, symbol string) (float64, error) {
	ctx = mc.getAuthContext(ctx)
	instruments, _, err := mc.client.InstrumentApi.InstrumentGetActive(ctx)
	if err != nil {
		return 0, err
	}
	for _, i := range instruments {
		if i.Symbol == symbol {
			return float64(i.MarkPrice), nil
		}
	}
	return 0, fmt.Errorf("could not find instrument")
}

func (mc *BitmexClient) fetchInstrument(ctx context.Context, symbol string) (*swagger.Instrument, error) {
	ctx = mc.getAuthContext(ctx)
	instruments, _, err := mc.client.InstrumentApi.InstrumentGetActive(ctx)
	if err != nil {
		return nil, err
	}
	for _, i := range instruments {
		if i.Symbol == symbol {
			return &i, nil
		}
	}
	return nil, fmt.Errorf("could not find instrument")
}

func (mc *BitmexClient) FetchDescriber(ctx context.Context, t *models.Trade) (*models.Describer, error) {
	errChan := make(chan error, 2)
	candleChan := make(chan *Candle, 1)
	symbolChan := make(chan *swagger.Instrument, 1)

	go func() {
		candle, err := mc.GetLastCompletedCandle(t.Symbol, t.Timeframe)
		if err != nil {
			errChan <- err
		}
		candleChan <- candle
	}()
	go func() {
		symbol, err := mc.GetSymbol(t.Symbol)
		if err != nil {
			errChan <- err
		}
		symbolChan <- symbol
	}()

	var candle *Candle
	var symbol *swagger.Instrument

	for i := 0; i < 2; i++ {
		select {
		case err := <-errChan:
			return nil, err
		case candle = <-candleChan:
		case symbol = <-symbolChan:
		}
	}

	timeframeDur, err := models.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}
	sp, err := t.CalculateStopPrice(candle.High, candle.Low)
	if err != nil {
		return nil, err
	}
	sl, err := t.CalculateStopLossPrice(candle.High, candle.Low, sp, false)
	if err != nil {
		return nil, err
	}
	tp, err := t.CalculateTakeProfitPrice(candle.High, candle.Low, sp, false)
	if err != nil {
		return nil, err
	}

	rsl, err := t.CalculateStopLossPrice(candle.High, candle.Low, sl, true)
	if err != nil {
		return nil, err
	}
	rtp, err := t.CalculateTakeProfitPrice(candle.High, candle.Low, sl, true)
	if err != nil {
		return nil, err
	}

	return &models.Describer{
		TradeID:           t.ID,
		Symbol:            t.Symbol,
		Size:              t.Size,
		Side:              t.Side,
		TakeProfitSize:    t.TakeProfitSize,
		StopLossSize:      t.StopLossSize,
		ReverseMultiplier: t.ReverseMultiplier,
		TimeFrameDur:      timeframeDur,

		OpenTime:               candle.OpenTime,
		CloseTime:              candle.CloseTime,
		Open:                   candle.Open,
		Close:                  candle.Close,
		High:                   candle.High,
		Low:                    candle.Low,
		StopPrice:              sp,
		TakeProfitPrice:        tp,
		StopLossPrice:          sl,
		ReverseStopLossPrice:   rsl,
		ReverseTakeProfitPrice: rtp,
		TickSize:               symbol.TickSize,
		LotSize:                float64(symbol.LotSize),
		MaxOrderQty:            float64(symbol.MaxOrderQty),
	}, nil
}
