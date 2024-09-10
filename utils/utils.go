package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func ConvertBinanceOrderID(orderID int64) string {
	return strconv.FormatInt(orderID, 10)
}

func ConvertOrderIDtoBinanceOrderID(orderID string) (int64, error) {
	return strconv.ParseInt(orderID, 10, 64)
}

func ConvertTime(timestamp int64) time.Time {
	// Convert milliseconds to seconds and nanoseconds
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1e6

	// Convert to time.Time
	return time.Unix(seconds, nanoseconds)
}

func FormatTimestamp(timestamp int64) string {
	t := time.Unix(0, timestamp*int64(time.Millisecond))

	formattedTime := t.Format("2000-01-01 15:04:05")

	return formattedTime
}

func FormatTime(t time.Time) string {
	return t.Format("2000-01-01 15:04:05")
}

func FriendlyDuration(duration time.Duration) string {
	// Convert to hours, minutes, seconds, etc.
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := int(duration.Milliseconds()) % 1000

	// Build a friendly string representation
	var friendlyDuration string
	if hours > 0 {
		friendlyDuration += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		friendlyDuration += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 {
		friendlyDuration += fmt.Sprintf("%ds ", seconds)
	}
	if milliseconds > 0 {
		friendlyDuration += fmt.Sprintf("%dms", milliseconds)
	}

	return friendlyDuration
}

func PrintStructFields(s interface{}) {
	val := reflect.ValueOf(s)
	// check if the passed value is a pointer
	if val.Kind() == reflect.Ptr {
		// Dereference the pointer to get the underlying struct
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		// Get the field name
		fieldName := typ.Field(i).Name

		// Get the field type
		fieldType := typ.Field(i).Type

		// Get the field value
		fieldValue := val.Field(i)

		fmt.Printf("%s (%s) = %v\n", fieldName, fieldType, fieldValue)
	}
}

func ParseTrade(tradeArgs []string) (*models.Trade, error) {
	var t models.Trade
	if len(tradeArgs) < 9 {
		return nil, fmt.Errorf("insufficient arguments provided; please ensure you have 9 parameters")
	}
	// account it should be string
	// m for bitmex
	// b for binance
	part1 := strings.TrimSpace(strings.ToLower(tradeArgs[0]))
	if len(part1) > 1 || (part1 != "m" && part1 != "b") {
		return nil, fmt.Errorf("invalid account value; use 'm' for BitMEX or 'b' for Binance")
	} else if part1 == "m" {
		t.Account = types.ACCOUNT_M
	} else if part1 == "b" {
		t.Account = types.ACCOUNT_B
	} else {
		// should never happen
		return nil, fmt.Errorf("unexpected internal error")
	}
	// pair
	part2 := strings.TrimSpace(strings.ToUpper(tradeArgs[1]))
	t.Symbol = part2
	// side
	part3 := strings.TrimSpace(strings.ToUpper(tradeArgs[2]))
	if part3 != types.SIDE_L && part3 != types.SIDE_S {
		return nil, fmt.Errorf("invalid side value; please enter 'long' or 'short'")
	} else {
		t.Side = part3
	}
	// candle
	part4 := strings.TrimSpace(tradeArgs[3])
	if !types.IsValidCandle(part4) {
		return nil, fmt.Errorf("invalid timeframe; valid values are: %s", types.GetValidCandlesString())
	} else {
		t.Timeframe = part4
	}
	// offset
	part5 := strings.TrimSpace(tradeArgs[4])
	offset, err := strconv.ParseFloat(part5, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid offset_entry; please provide a numeric value")
	} else {
		t.Offset = offset
	}
	// size percent
	part6 := strings.TrimSpace(tradeArgs[5])
	size_percent, err := strconv.Atoi(part6)
	if err != nil {
		return nil, fmt.Errorf("invalid size; please provide a percentage value (e.g., 5)")
	} else if size_percent <= 0 || size_percent > 50 {
		return nil, fmt.Errorf("invalid size; please provide a value between 1 and 50")
	} else {
		t.Size = size_percent
	}

	// stop-loss percent
	part7 := strings.TrimSpace(tradeArgs[6])
	stop_percent, err := strconv.Atoi(part7)
	if err != nil {
		return nil, fmt.Errorf("invalid stop-loss percent; please provide a numeric value (e.g., 105)")
	} else if stop_percent < 100 {
		return nil, fmt.Errorf("invalid stop-loss percent; must be 100 or greater")
	} else {
		t.StopLoss = stop_percent
	}

	// target-point percent
	part8 := strings.TrimSpace(tradeArgs[7])
	target_percent, err := strconv.Atoi(part8)
	if err != nil {
		return nil, fmt.Errorf("invalid target-point percent; please provide a numeric value (e.g., 105)")
	} else if target_percent < 100 {
		return nil, fmt.Errorf("invalid target-point percent; must be 100 or greater")
	} else {
		t.TakeProfit = target_percent
	}

	// reverse-multiplier
	part9 := strings.TrimSpace(tradeArgs[8])
	reverse_multiplier, err := strconv.Atoi(part9)
	if err != nil {
		return nil, fmt.Errorf("invalid reverse_multiplier; please provide a value of 1 or 2")
	} else if reverse_multiplier <= 0 || reverse_multiplier > 2 {
		return nil, fmt.Errorf("invalid reverse_multiplier; must be 1 or 2")
	} else {
		t.ReverseMultiplier = reverse_multiplier
	}

	t.State = types.STATE_IDLE

	return &t, nil
}
