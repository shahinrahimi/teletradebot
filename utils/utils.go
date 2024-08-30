package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
)

func ConvertTime(timestamp int64) time.Time {
	// Convert milliseconds to seconds and nanoseconds
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1e6

	// Convert to time.Time
	return time.Unix(seconds, nanoseconds)
}

func PrintStructFields(s interface{}) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

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
		return nil, fmt.Errorf("the length of args is not sufficient for parsing")
	}
	// account it should be string
	// m for bitmex
	// b for binance
	part1 := strings.TrimSpace(strings.ToLower(tradeArgs[0]))
	if len(part1) > 1 || (part1 != "m" && part1 != "b") {
		return nil, fmt.Errorf("the valid value for account should be 'm' => bitmex, 'b' => binance")
	} else if part1 == "m" {
		t.Account = types.ACCOUNT_M
	} else if part1 == "b" {
		t.Account = types.ACCOUNT_B
	} else {
		// should never happen
		return nil, fmt.Errorf("internal error")
	}
	// pair
	// TODO maybe add check if pair exist on the tickers
	part2 := strings.TrimSpace(strings.ToUpper(tradeArgs[1]))
	t.Pair = part2
	// side
	part3 := strings.TrimSpace(strings.ToUpper(tradeArgs[2]))
	if part3 != types.SIDE_L && part3 != types.SIDE_S {
		return nil, fmt.Errorf("the valid value for side should be 'long' or 'short'")
	} else {
		t.Side = part3
	}
	// candle
	part4 := strings.TrimSpace(tradeArgs[3])
	if !types.IsValidCandle(part4) {
		return nil, fmt.Errorf("the valid value for candle should be %s", types.GetValidCandlesString())
	} else {
		t.Candle = part4
	}
	// offset
	part5 := strings.TrimSpace(tradeArgs[4])
	offset, err := strconv.ParseFloat(part5, 64)
	if err != nil {
		return nil, fmt.Errorf("the valid value for offset_entry should be amount (float or integer)")
	} else {
		t.Offset = float32(offset)
	}
	// size percent
	part6 := strings.TrimSpace(tradeArgs[5])
	size_percent, err := strconv.Atoi(part6)
	if err != nil {
		return nil, fmt.Errorf("the valid value for size should be amount in percent (e.g 5)")
	} else if size_percent <= 0 || size_percent > 50 {
		return nil, fmt.Errorf("the valid value for size should be a non-zero none-negative number (max: 50)")
	} else {
		t.SizePercent = size_percent
	}

	// stop-loss percent
	part7 := strings.TrimSpace(tradeArgs[6])
	stop_percent, err := strconv.Atoi(part7)
	if err != nil {
		return nil, fmt.Errorf("the valid value for stop-loss percent should be amount in percent (e.g 105)")
	} else if stop_percent < 100 {
		return nil, fmt.Errorf("the valid value for stop-loss percent should be a non-zero none-negative number (min: 100)")
	} else {
		t.SLPercent = stop_percent
	}

	// target-point percent
	part8 := strings.TrimSpace(tradeArgs[7])
	target_percent, err := strconv.Atoi(part8)
	if err != nil {
		return nil, fmt.Errorf("the valid value for target-point percent should be amount in percent (e.g 105)")
	} else if target_percent < 100 {
		return nil, fmt.Errorf("the valid value for target-point percent should be a non-zero none-negative number (min: 100)")
	} else {
		t.TPPercent = target_percent
	}

	// reverse-multiplier
	part9 := strings.TrimSpace(tradeArgs[8])
	reverse_multiplier, err := strconv.Atoi(part9)
	if err != nil {
		return nil, fmt.Errorf("the valid value for reverse_multiplier should be number (1 or 2)")
	} else if reverse_multiplier <= 0 || reverse_multiplier > 2 {
		return nil, fmt.Errorf("the valid value for target-point percent should be a non-zero none-negative number (1 or 2)")
	} else {
		t.ReverseMultiplier = reverse_multiplier
	}

	t.State = types.STATE_IDLE

	return &t, nil
}
