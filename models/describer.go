package models

import (
	"fmt"
	"math"
	"time"

	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

type Describer struct {
	// these fields copy from trade
	TradeID           int64
	Symbol            string
	Side              string
	Size              int
	StopLossSize      int
	TakeProfitSize    int
	ReverseMultiplier int
	TimeFrameDur      time.Duration

	// these fields fetched from exchange
	OpenTime  time.Time
	CloseTime time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	// these fields calculated
	StopPrice              float64
	StopLossPrice          float64
	TakeProfitPrice        float64
	ReverseStopLossPrice   float64
	ReverseTakeProfitPrice float64

	PricePrecision    int // use in binance exchange
	QuantityPrecision int // use in binance exchange

	TickSize    float64 // use in bitmex exchange
	LotSize     float64 // use in bitmex exchange
	MaxOrderQty float64 // use in bitmex exchange
}

func (d *Describer) AdjustPriceForBinance(value float64) string {
	return fmt.Sprintf("%.*f", d.PricePrecision, value)
}

func (d *Describer) AdjustQuantityForBinance(value float64) string {
	return fmt.Sprintf("%.*f", d.QuantityPrecision, value)
}

func (d *Describer) AdjustPriceForBitmex(value float64) float64 {
	return math.Round(value/d.TickSize) * d.TickSize
}

func (d *Describer) AdjustQuantityForBitmex(value float64) float64 {
	return math.Floor(value/d.LotSize) * d.LotSize
}

func (d *Describer) getPriceString(price float64) string {
	var priceStr string
	if d.PricePrecision > 0 {
		priceStr = fmt.Sprintf("%.*f", d.PricePrecision, price)
	} else if d.TickSize > 0 {
		p := math.Floor(price/d.TickSize) * d.TickSize
		precision := int(math.Abs(math.Log10(d.TickSize)))
		priceStr = fmt.Sprintf("%.*f", precision, p)
	} else {
		priceStr = fmt.Sprintf("%f", price)
	}
	return priceStr
}

func (d *Describer) CalculateExpiration() time.Duration {
	return d.TimeFrameDur + time.Until(d.CloseTime)
}
func (d *Describer) ToString() string {

	format := "2006-01-02 15:04:05"
	from := d.OpenTime.Local().Format(format)
	till := d.CloseTime.Local().Format(format)

	size := fmt.Sprintf("%.1f%%", float64(d.Size))
	reverseSize := fmt.Sprintf("%.1f%%", float64(d.ReverseMultiplier)*float64(d.Size))
	stopLossSize := fmt.Sprintf("%.1f%%", float64((d.StopLossSize - 100)))
	takeProfitSize := fmt.Sprintf("%.1f%%", float64((d.TakeProfitSize - 100)))

	open := d.getPriceString(d.Open)
	close := d.getPriceString(d.Close)
	high := d.getPriceString(d.High)
	low := d.getPriceString(d.Low)

	stopPrice := d.getPriceString(d.StopPrice)
	stopLossPrice := d.getPriceString(d.StopLossPrice)
	takeProfitPrice := d.getPriceString(d.TakeProfitPrice)
	reverseStopLossPrice := d.getPriceString(d.ReverseStopLossPrice)
	reverseTakeProfitPrice := d.getPriceString(d.ReverseTakeProfitPrice)

	var reverseSide string
	if d.Side == types.SIDE_L {
		reverseSide = types.SIDE_S
	} else {
		reverseSide = types.SIDE_L
	}

	var expiration string
	if d.CalculateExpiration() > 0 {
		expiration = utils.FriendlyDuration(d.CalculateExpiration())
	} else {
		expiration = "∞"
	}

	msg := fmt.Sprintf("Trade ID %d\n\n", d.TradeID)
	msg = fmt.Sprintf("%sFrom:  %s\nTill:  %s\nOpen:  %s\nHigh:  %s\nLow:  %s\nClose:  %s\n\n", msg, from, till, open, high, low, close)
	msg = fmt.Sprintf("%sTrading:\n", msg)
	msg = fmt.Sprintf("%sEntry %s at %s with %s of balance.\n", msg, d.Side, stopPrice, size)
	msg = fmt.Sprintf("%sTP at %s with %s.\n", msg, takeProfitPrice, takeProfitSize)
	msg = fmt.Sprintf("%sSL at %s with %s.\n\n", msg, stopLossPrice, stopLossSize)
	if d.ReverseMultiplier > 0 {
		msg = fmt.Sprintf("%sReverse:\n", msg)
		msg = fmt.Sprintf("%sEntry %s at %s with %s of balance.\n", msg, reverseSide, stopLossPrice, reverseSize)
		msg = fmt.Sprintf("%sTP at %s with %s.\n", msg, reverseTakeProfitPrice, takeProfitSize)
		msg = fmt.Sprintf("%sSL at %s with %s.\n\n", msg, reverseStopLossPrice, stopLossSize)
	}
	msg = fmt.Sprintf("%sExpiration: %s", msg, expiration)
	return msg
}
