package models

import (
	"fmt"
	"time"
)

type Describer struct {
	OpenTime          time.Time
	CloseTime         time.Time
	Open              float64
	High              float64
	Low               float64
	Close             float64
	StopPrice         float64
	StopLossPrice     float64
	TakeProfitPrice   float64
	CandleDuration    time.Duration
	PricePrecision    int     // use in binance exchange
	QuantityPrecision int     // use in binance exchange
	TickSize          float64 // use in bitmex exchange
}

func (d *Describer) CalculateExpiration() time.Duration {
	return d.CandleDuration + time.Until(d.CloseTime)
}

func (d *Describer) ToString(t *Trade) string {

	format := "2006-01-02 15:04:05"
	from := d.OpenTime.Local().Format(format)
	till := d.CloseTime.Local().Format(format)

	size := fmt.Sprintf("%.1f%%", float64(t.Size))
	stopLossSize := fmt.Sprintf("%.1f%%", float64((t.StopLoss - 100)))
	takeProfitSize := fmt.Sprintf("%.1f%%", float64((t.TakeProfit - 100)))

	open := fmt.Sprintf("%f", d.Open)
	close := fmt.Sprintf("%f", d.Close)
	high := fmt.Sprintf("%f", d.High)
	low := fmt.Sprintf("%f", d.Low)

	stopPrice := fmt.Sprintf("%f", d.StopPrice)
	stopLossPrice := fmt.Sprintf("%f", d.StopLossPrice)
	takeProfitPrice := fmt.Sprintf("%f", d.TakeProfitPrice)

	msg := fmt.Sprintf("Trade ID %d\n\n", t.ID)
	msg = fmt.Sprintf("%sFrom:  %s\nTill:  %s\nOpen:  %s\nHigh:  %s\nLow:  %s\nClose:  %s\n\n", msg, from, till, open, high, low, close)
	msg = fmt.Sprintf("%sTrading:\n", msg)
	msg = fmt.Sprintf("%sEntry %s at %s with %s of balance.\n", msg, t.Side, stopPrice, size)
	msg = fmt.Sprintf("%sTP at %s with %s.\n", msg, takeProfitPrice, takeProfitSize)
	msg = fmt.Sprintf("%sSL at %s with %s.\n", msg, stopLossPrice, stopLossSize)
	return msg
}
