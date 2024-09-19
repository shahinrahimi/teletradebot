package bitmex

import (
	"fmt"
	"sync"
	"time"

	"github.com/shahinrahimi/teletradebot/utils"
)

type Candle struct {
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	OpenTime  time.Time
	CloseTime time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Timeframe time.Duration
	Completed bool
}
type SymbolCandle struct {
	Symbol     string
	timeframes map[time.Duration][]*Candle // last candle is current candle

}

var (
	symbolCandles = make(map[string]SymbolCandle)
	timeframes    = []time.Duration{
		time.Minute,
		time.Minute * 3,
		time.Minute * 5,
		time.Minute * 15,
		time.Minute * 30,
		time.Hour,
	}
	mu sync.RWMutex
)

func getCurrentCandle(candles []*Candle, timeframe time.Duration, timestamp time.Time) *Candle {
	if len(candles) > 0 {
		targetCandle := candles[len(candles)-1]
		if targetCandle.OpenTime.Truncate(timeframe) == timestamp.Truncate(timeframe) {
			return candles[len(candles)-1]
		} else {
			targetCandle.Completed = true
		}
	}
	// create a first candle for the symbol
	return &Candle{
		Open:      0,
		OpenTime:  timestamp.Truncate(timeframe),
		CloseTime: timestamp.Add(timeframe).Truncate(timeframe),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (mc *BitmexClient) UpdateCandles(symbol string, markPrice float64, timestamp time.Time) {
	mu.Lock()
	defer mu.Unlock()

	sc, ok := symbolCandles[symbol]
	if !ok {
		sc = SymbolCandle{
			Symbol:     symbol,
			timeframes: make(map[time.Duration][]*Candle),
		}
		for _, timeframe := range timeframes {
			sc.timeframes[timeframe] = make([]*Candle, 0)
		}
		symbolCandles[symbol] = sc
	}

	sc = symbolCandles[symbol]
	for _, timeframe := range timeframes {
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
	}
}

func (mc *BitmexClient) GetLastClosedCandle(symbol string, timeframe time.Duration) (*Candle, error) {
	sc, ok := symbolCandles[symbol]
	if !ok {
		return nil, fmt.Errorf("symbol %s not found", symbol)
	}
	if len(sc.timeframes[timeframe]) <= 2 {
		return nil, fmt.Errorf("symbol %s with timeframe %s is not a complete timeframe", symbol, utils.FriendlyDuration(timeframe))
	}
	return sc.timeframes[timeframe][len(sc.timeframes[timeframe])-2], nil
}
