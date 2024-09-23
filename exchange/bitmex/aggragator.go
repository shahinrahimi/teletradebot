package bitmex

import (
	"fmt"
	"sync"
	"time"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

type Candle struct {
	Open        float64
	High        float64
	Low         float64
	Close       float64
	Volume      float64
	OpenTime    time.Time
	CloseTime   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Timeframe   time.Duration
	Completed   bool
	FirstCandle bool
}
type SymbolCandle struct {
	Symbol     string
	timeframes map[models.TimeframeType][]*Candle // last candle is current candle

}

var (
	symbolCandles = make(map[string]SymbolCandle)
	mu            sync.RWMutex
)

func getCurrentCandle(candles []*Candle, timeframe models.TimeframeType, timestamp time.Time) *Candle {
	timeframeDur, err := models.GetDuration(timeframe)
	if err != nil {
		panic(err)
	}
	if len(candles) > 0 {
		targetCandle := candles[len(candles)-1]
		if targetCandle.OpenTime.Truncate(timeframeDur) == timestamp.Truncate(timeframeDur) {
			return targetCandle
		} else {
			targetCandle.Completed = true
		}
	}
	var firstCandle = false
	if len(candles) == 0 {
		firstCandle = true
	}
	// create a first candle for the symbol
	return &Candle{
		Open:        0,
		OpenTime:    timestamp.Truncate(timeframeDur),
		CloseTime:   timestamp.Add(timeframeDur).Truncate(timeframeDur),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		FirstCandle: firstCandle,
	}
}

func (mc *BitmexClient) UpdateCandles(symbol string, markPrice float64, timestamp time.Time) {
	mu.Lock()
	defer mu.Unlock()

	sc, ok := symbolCandles[symbol]
	if !ok {
		sc = SymbolCandle{
			Symbol:     symbol,
			timeframes: make(map[models.TimeframeType][]*Candle),
		}
		for _, timeframe := range models.GetValidTimeframes() {
			sc.timeframes[timeframe] = make([]*Candle, 0)
		}
		symbolCandles[symbol] = sc
	}

	sc = symbolCandles[symbol]
	for _, timeframe := range models.GetValidTimeframes() {
		candles := sc.timeframes[timeframe]
		currentCandle := getCurrentCandle(candles, timeframe, timestamp)
		// close price actually is the live price
		currentCandle.Close = markPrice
		// if the candle is new, initialize its open high, low, close
		if currentCandle.Open == 0 {
			sc.timeframes[timeframe] = append(sc.timeframes[timeframe], currentCandle)
			currentCandle.Open = markPrice
			currentCandle.High = markPrice
			currentCandle.Low = markPrice
		}
		if markPrice > currentCandle.High {
			currentCandle.High = markPrice
		}
		if markPrice < currentCandle.Low {
			currentCandle.Low = markPrice
		}
		currentCandle.UpdatedAt = time.Now().UTC()

		// Keep only the last 2 candles (current and previous)
		if len(sc.timeframes[timeframe]) > 2 {
			sc.timeframes[timeframe] = sc.timeframes[timeframe][1:]
		}
	}
}

func (mc *BitmexClient) GetLastClosedCandle(symbol string, timeframe models.TimeframeType) (*Candle, error) {
	mu.RLock()
	defer mu.RUnlock()
	timeframeDur, err := models.GetDuration(timeframe)
	if err != nil {
		return nil, err
	}
	sc, ok := symbolCandles[symbol]
	if !ok {
		return nil, fmt.Errorf("symbol %s not found", symbol)
	}
	for _, c := range sc.timeframes[timeframe] {
		if c.Completed && !c.FirstCandle {
			return c, nil
		}
	}
	switch len(sc.timeframes[timeframe]) {
	case 1:
		dur := -time.Since(sc.timeframes[timeframe][0].CloseTime) + timeframeDur
		return nil, &types.BotError{
			Msg: fmt.Sprintf("Symbol %s with timeframe %s has not completed yet. Please wait for %s.", symbol, timeframe, utils.FriendlyDuration(dur)),
		}
	case 2:
		dur := -time.Since(sc.timeframes[timeframe][1].CloseTime)
		return nil, &types.BotError{
			Msg: fmt.Sprintf("Symbol %s with timeframe %s has not completed yet. Please wait for %s.", symbol, timeframe, utils.FriendlyDuration(dur)),
		}
	default:
		return nil, &types.BotError{
			Msg: fmt.Sprintf("Symbol %s with timeframe %s has not completed yet.", symbol, timeframe),
		}
	}
}
