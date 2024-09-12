package types

import (
	"fmt"
	"time"

	"github.com/shahinrahimi/teletradebot/models"
)

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
	STATE_CANCELED  string = `canceled`
	STATE_PLACED    string = `placed`
	STATE_FILLED    string = `filled`
	STATE_REVERTING string = `reverting`
	STATE_STOPPED   string = `stopped`
	STATE_PROFITED  string = `profited`
)

var ExpireDuration = map[string]time.Duration{
	CANDLE_1MIN:  time.Minute * 1,
	CANDLE_3MIN:  time.Minute * 3,
	CANDLE_5MIN:  time.Minute * 5,
	CANDLE_15MIN: time.Minute * 15,
	CANDLE_30MIN: time.Minute * 30,
	CANDLE_1H:    time.Hour * 1,
	CANDLE_2H:    time.Hour * 2,
	CANDLE_4H:    time.Hour * 4,
	CANDLE_6H:    time.Hour * 6,
	CANDLE_8H:    time.Hour * 8,
	CANDLE_12H:   time.Hour * 12,
	CANDLE_1D:    time.Hour * 24,
	CANDLE_3D:    time.Hour * 24 * 3,
	CANDLE_1W:    time.Hour * 24 * 7,
	CANDLE_1M:    time.Hour * 24 * 30, // Approximation, as months vary in length
}

func GetDuration(candle string) (time.Duration, error) {
	if duration, exists := ExpireDuration[candle]; exists {
		return duration, nil
	} else {
		return 0, fmt.Errorf("Candle interval %s not found", candle)
	}
}

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

type BotMessage struct {
	ChatID int64
	MsgStr string
}

type TradeDescriber struct {
	From  time.Time
	Till  time.Time
	Open  string
	High  string
	Low   string
	Close string
	SP    string // strop price or entry
	TP    string // take-profit price
	SL    string // take-loss price
}

var (
	TradeDescribers = map[int64]*TradeDescriber{}
)

func (td *TradeDescriber) ToTelegramString(t *models.Trade) string {
	sizeStr := fmt.Sprintf("%.1f%%", float64(t.Size))
	slStr := fmt.Sprintf("%.1f%%", float64((t.StopLoss - 100)))
	tpStr := fmt.Sprintf("%.1f%%", float64((t.TakeProfit - 100)))

	format := "2006-01-02 15:04:05"
	FromStr := td.From.Local().Format(format)
	TillStr := td.Till.Local().Format(format)

	msg := fmt.Sprintf("Trade ID %d\n\n", t.ID)
	msg = fmt.Sprintf("%s From:  %s\n Till:  %s\n Open:  %s\n High:  %s\n Low:  %s\n Close:  %s\n\n", msg, FromStr, TillStr, td.Open, td.High, td.Low, td.Close)
	msg = fmt.Sprintf("%sTrading:\n", msg)
	msg = fmt.Sprintf("%s Entry %s at %s with %s of balance.\n", msg, t.Side, td.SP, sizeStr)
	msg = fmt.Sprintf("%s TP at %s with %s.\n", msg, td.TP, tpStr)
	msg = fmt.Sprintf("%s SL at %s with %s.\n", msg, td.SL, slStr)
	return msg
}
