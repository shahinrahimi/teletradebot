package models

import (
	"fmt"
	"math"
	"time"
)

type Describer struct {
	OpenTime        time.Time
	CloseTime       time.Time
	Open            float64
	High            float64
	Low             float64
	Close           float64
	StopPrice       float64
	StopLossPrice   float64
	TakeProfitPrice float64
	CandleDuration  time.Duration

	PricePrecision    int // use in binance exchange
	QuantityPrecision int // use in binance exchange

	TickSize float64 // use in bitmex exchange
	LotSize  float64 // use in bitmex exchange
}

func (d *Describer) CalculateExpiration() time.Duration {
	return d.CandleDuration + time.Until(d.CloseTime)
}

func (d *Describer) getPriceString(price float64) string {
	var priceStr string
	if d.PricePrecision > 0 {
		pricePrecision := math.Pow10(int(-d.PricePrecision))
		p := math.Floor(price/pricePrecision) * pricePrecision
		priceStr = fmt.Sprintf("%.*f", d.PricePrecision, p)
	} else {
		priceStr = fmt.Sprintf("%f", price)
	}
	return priceStr
}

func (d *Describer) ToString(t *Trade) string {

	format := "2006-01-02 15:04:05"
	from := d.OpenTime.Local().Format(format)
	till := d.CloseTime.Local().Format(format)

	size := fmt.Sprintf("%.1f%%", float64(t.Size))
	stopLossSize := fmt.Sprintf("%.1f%%", float64((t.StopLoss - 100)))
	takeProfitSize := fmt.Sprintf("%.1f%%", float64((t.TakeProfit - 100)))

	open := d.getPriceString(d.Open)
	close := d.getPriceString(d.Close)
	high := d.getPriceString(d.High)
	low := d.getPriceString(d.Low)

	stopPrice := d.getPriceString(d.StopPrice)
	stopLossPrice := d.getPriceString(d.StopLossPrice)
	takeProfitPrice := d.getPriceString(d.TakeProfitPrice)

	var expiration string
	if d.CalculateExpiration() > 0 {
		expiration = friendlyDuration(d.CalculateExpiration())
	} else {
		expiration = "∞"
	}

	msg := fmt.Sprintf("Trade ID %d\n\n", t.ID)
	msg = fmt.Sprintf("%sFrom:  %s\nTill:  %s\nOpen:  %s\nHigh:  %s\nLow:  %s\nClose:  %s\n\n", msg, from, till, open, high, low, close)
	msg = fmt.Sprintf("%sTrading:\n", msg)
	msg = fmt.Sprintf("%sEntry %s at %s with %s of balance.\n", msg, t.Side, stopPrice, size)
	msg = fmt.Sprintf("%sTP at %s with %s.\n", msg, takeProfitPrice, takeProfitSize)
	msg = fmt.Sprintf("%sSL at %s with %s.\n\n", msg, stopLossPrice, stopLossSize)
	msg = fmt.Sprintf("%sExpiration: %s", msg, expiration)
	return msg
}

func friendlyDuration(duration time.Duration) string {
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
