package types

const (
	ACCOUNT_B string = `Binance`
	ACCOUNT_M string = `Bitmex`

	SIDE_L string = `LONG`
	SIDE_S string = `SHORT`

	CANDLE_1MIN  string = `1m`
	CANDLE_3MIN  string = `3m`
	CANDLE_5MIN  string = `5m`
	CANDLE_15MIN string = `15m`
	CANDLE_30MIN string = `30m`
	CANDLE_1H    string = `1h`
	CANDLE_2H    string = `2h`
	CANDLE_4H    string = `4h`
	CANDLE_6H    string = `6h`
	CANDLE_8H    string = `8h`
	CANDLE_12H   string = `12h`
	CANDLE_1D    string = `1d`
	CANDLE_3D    string = `3d`
	CANDLE_1W    string = `1w`
	CANDLE_1M    string = `1M`

	STATE_IDLE      string = `idle`
	STATE_PLACED    string = `placed`
	STATE_FILLED    string = `filled`
	STATE_REVERTING string = `reverting`
)

func GetValidCandles() []string {
	validCandles := []string{
		CANDLE_1MIN,
		CANDLE_3MIN,
		CANDLE_5MIN,
		CANDLE_15MIN,
		CANDLE_30MIN,
		CANDLE_1H,
		CANDLE_2H,
		CANDLE_4H,
		CANDLE_6H,
		CANDLE_8H,
		CANDLE_12H,
		CANDLE_1D,
		CANDLE_3D,
		CANDLE_1W,
		CANDLE_1M,
	}
	return validCandles
}

func GetValidCandlesString() string {
	var candleStr string = ""
	var plusStr string
	validCandles := GetValidCandles()
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

func IsValidCandle(candle string) bool {
	for _, c := range GetValidCandles() {
		if candle == c {
			return true
		}
	}
	return false
}
