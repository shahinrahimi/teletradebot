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

func (mc *BitmexClient) GetLastCompletedCandle(symbol string, timeframe models.TimeframeType) (*Candle, error) {
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
