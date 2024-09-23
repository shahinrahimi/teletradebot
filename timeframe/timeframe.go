package timeframe

import (
	"fmt"
	"time"
)

const (
	TIMEFRAME_1MIN  string = `1m`
	TIMEFRAME_3MIN  string = `3m`
	TIMEFRAME_5MIN  string = `5m`
	TIMEFRAME_15MIN string = `15m`
	TIMEFRAME_30MIN string = `30m`
	TIMEFRAME_1H    string = `1h`
	TIMEFRAME_2H    string = `2h`
	TIMEFRAME_4H    string = `4h`
	TIMEFRAME_6H    string = `6h`
	TIMEFRAME_8H    string = `8h`
	TIMEFRAME_12H   string = `12h`
	TIMEFRAME_1D    string = `1d`
	TIMEFRAME_3D    string = `3d`
	TIMEFRAME_1W    string = `1w`
	TIMEFRAME_1M    string = `1M`
)

var Timeframes = map[string]time.Duration{
	TIMEFRAME_1MIN:  time.Minute * 1,
	TIMEFRAME_3MIN:  time.Minute * 3,
	TIMEFRAME_5MIN:  time.Minute * 5,
	TIMEFRAME_15MIN: time.Minute * 15,
	TIMEFRAME_30MIN: time.Minute * 30,
	TIMEFRAME_1H:    time.Hour * 1,
	TIMEFRAME_2H:    time.Hour * 2,
	TIMEFRAME_4H:    time.Hour * 4,
	TIMEFRAME_6H:    time.Hour * 6,
	TIMEFRAME_8H:    time.Hour * 8,
	TIMEFRAME_12H:   time.Hour * 12,
	TIMEFRAME_1D:    time.Hour * 24,
	TIMEFRAME_3D:    time.Hour * 24 * 3,
	TIMEFRAME_1W:    time.Hour * 24 * 7,
	TIMEFRAME_1M:    time.Hour * 24 * 30, // Approximation, as months vary in length
}

func GetDuration(timeframe string) (time.Duration, error) {
	if duration, exists := Timeframes[timeframe]; exists {
		return duration, nil
	} else {
		return 0, fmt.Errorf("timeframe %s not found", timeframe)
	}
}

func getValidTimeframes() []string {
	validCandles := []string{
		TIMEFRAME_1MIN,
		TIMEFRAME_3MIN,
		TIMEFRAME_5MIN,
		TIMEFRAME_15MIN,
		TIMEFRAME_30MIN,
		TIMEFRAME_1H,
		TIMEFRAME_2H,
		TIMEFRAME_4H,
		TIMEFRAME_6H,
		TIMEFRAME_8H,
		TIMEFRAME_12H,
		TIMEFRAME_1D,
		TIMEFRAME_3D,
		TIMEFRAME_1W,
		TIMEFRAME_1M,
	}
	return validCandles
}

func GetValidTimeframesString() string {
	var candleStr string = ""
	var plusStr string
	validCandles := getValidTimeframes()
	for index, c := range validCandles {
		if index == len(validCandles)-1 {
			plusStr = ", "
		} else {
			plusStr = " and "
		}
		candleStr = candleStr + plusStr + c
	}
	return candleStr
}

func IsValidTimeframe(candle string) bool {
	for _, c := range getValidTimeframes() {
		if candle == c {
			return true
		}
	}
	return false
}
