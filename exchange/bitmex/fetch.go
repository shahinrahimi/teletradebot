package bitmex

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
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

func (mc *BitmexClient) FetchInterpreter(ctx context.Context, t *models.Trade) (*models.Interpreter, error) {
	errChan := make(chan error, 4)
	balanceChan := make(chan float64, 1)
	priceChan := make(chan float64, 1)
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
	go func() {
		balance, err := mc.fetchBalance(ctx)
		if err != nil {
			errChan <- err
		}
		balanceChan <- balance
	}()
	go func() {
		price, err := mc.fetchPrice(ctx, t.Symbol)
		if err != nil {
			errChan <- err
		}
		priceChan <- price
	}()

	var candle *Candle
	var symbol *swagger.Instrument
	var balance float64
	var price float64
	for i := 0; i < 4; i++ {
		select {
		case err := <-errChan:
			return nil, err
		case candle = <-candleChan:
		case symbol = <-symbolChan:
		case balance = <-balanceChan:
		case price = <-priceChan:
		}
	}

	timeframeDur, err := models.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}
	ep, err := t.CalculateEntryPrice(candle.High, candle.Low)
	if err != nil {
		return nil, err
	}
	sl, err := t.CalculateStopLossPrice(candle.High, candle.Low, ep, false)
	if err != nil {
		return nil, err
	}
	tp, err := t.CalculateTakeProfitPrice(candle.High, candle.Low, ep, false)
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

	contractSize, exist := config.ContractSizes[t.Symbol]
	if !exist {
		return nil, fmt.Errorf("contract size not found for symbol %s", t.Symbol)
	}

	balance = balance / 1000000 // balance in USDT

	size := balance * (float64(t.Size) / 100)
	quantity := size / (price * contractSize)
	rQuantity := quantity * float64(t.ReverseMultiplier)

	return &models.Interpreter{
		Balance:         balance,
		Price:           price,
		Quantity:        quantity,
		ReverseQuantity: rQuantity,
		Exchange:        types.ExchangeBitmex,

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
		EntryPrice:             ep,
		TakeProfitPrice:        tp,
		StopLossPrice:          sl,
		ReverseEntryPrice:      sl,
		ReverseStopLossPrice:   rsl,
		ReverseTakeProfitPrice: rtp,
		TickSize:               symbol.TickSize,
		LotSize:                float64(symbol.LotSize),
		MaxOrderQty:            float64(symbol.MaxOrderQty),
	}, nil
}
