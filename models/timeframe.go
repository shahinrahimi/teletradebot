package models

import (
	"fmt"
	"time"
)

type TimeframeType string

const (
	TIMEFRAME_1MIN  TimeframeType = `1m`
	TIMEFRAME_3MIN  TimeframeType = `3m`
	TIMEFRAME_5MIN  TimeframeType = `5m`
	TIMEFRAME_15MIN TimeframeType = `15m`
	TIMEFRAME_30MIN TimeframeType = `30m`
	TIMEFRAME_1H    TimeframeType = `1h`
	TIMEFRAME_2H    TimeframeType = `2h`
	TIMEFRAME_4H    TimeframeType = `4h`
	TIMEFRAME_6H    TimeframeType = `6h`
	TIMEFRAME_8H    TimeframeType = `8h`
	TIMEFRAME_12H   TimeframeType = `12h`
	TIMEFRAME_1D    TimeframeType = `1d`
	TIMEFRAME_3D    TimeframeType = `3d`
	TIMEFRAME_1W    TimeframeType = `1w`
	TIMEFRAME_1M    TimeframeType = `1M`
)

var Timeframes = map[TimeframeType]time.Duration{
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

func GetDuration(timeframe TimeframeType) (time.Duration, error) {
	if duration, exists := Timeframes[timeframe]; exists {
		return duration, nil
	} else {
		return 0, fmt.Errorf("timeframe %s not found", timeframe)
	}
}

func GetValidTimeframes() []TimeframeType {
	validCandles := []TimeframeType{
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
	var timeframeStr string = ""
	var plusStr string
	validTimeframes := GetValidTimeframes()
	for index, t := range validTimeframes {
		if index == len(validTimeframes)-1 {
			plusStr = ", "
		} else {
			plusStr = " and "
		}
		timeframeStr = timeframeStr + plusStr + string(t)
	}
	return timeframeStr
}

func IsValidTimeframe(timeframeStr string) bool {
	for _, c := range GetValidTimeframes() {
		if timeframeStr == string(c) {
			return true
		}
	}
	return false
}
