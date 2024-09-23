package models

import (
	"fmt"
	"math"
	"time"

	"github.com/shahinrahimi/teletradebot/utils"
)

type Describer struct {
	// these fields copy from trade
	TradeID        int64
	Symbol         string
	Side           string
	Size           int
	StopLossSize   int
	TakeProfitSize int
	CandleDuration time.Duration

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

	//PricePrecision
	//
	PricePrecision    float64 // use in binance exchange
	QuantityPrecision float64 // use in binance exchange

	//TickSize float64 // use in bitmex exchange
	//LotSize  float64 // use in bitmex exchange
}

func (d *Describer) CalculateExpiration() time.Duration {
	return d.CandleDuration + time.Until(d.CloseTime)
}

// func (d *Describer) getPriceString(price float64) string {
// 	var priceStr string
// 	if d.PricePrecision > 0 {
// 		p := math.Floor(price/d.PricePrecision) * d.PricePrecision
// 		precision := int(math.Abs(math.Log10(d.PricePrecision)))
// 		priceStr = fmt.Sprintf("%.*f", precision, p)
// 	} else {
// 		priceStr = fmt.Sprintf("%f", price)
// 	}
// 	return priceStr
// }

func (d *Describer) GetValueWithPricePrecision(value float64) float64 {
	return math.Floor(value/d.PricePrecision) * d.PricePrecision
}

func (d *Describer) GetValueWithPricePrecisionString(value float64) string {
	return fmt.Sprintf("%.*f", int(math.Abs(math.Log10(d.PricePrecision))), d.GetValueWithPricePrecision(value))
}

func (d *Describer) GetValueWithQuantityPrecision(value float64) float64 {
	return math.Floor(value/d.QuantityPrecision) * d.QuantityPrecision
}

func (d *Describer) GetValueWithQuantityPrecisionString(value float64) string {
	return fmt.Sprintf("%.*f", int(math.Abs(math.Log10(d.QuantityPrecision))), d.GetValueWithQuantityPrecision(value))
}
func (d *Describer) ToString() string {

	format := "2006-01-02 15:04:05"
	from := d.OpenTime.Local().Format(format)
	till := d.CloseTime.Local().Format(format)

	size := fmt.Sprintf("%.1f%%", float64(d.Size))
	stopLossSize := fmt.Sprintf("%.1f%%", float64((d.StopLossSize - 100)))
	takeProfitSize := fmt.Sprintf("%.1f%%", float64((d.TakeProfitSize - 100)))

	open := d.GetValueWithPricePrecisionString(d.Open)
	close := d.GetValueWithPricePrecisionString(d.Close)
	high := d.GetValueWithPricePrecisionString(d.High)
	low := d.GetValueWithPricePrecisionString(d.Low)

	stopPrice := d.GetValueWithPricePrecisionString(d.StopPrice)
	stopLossPrice := d.GetValueWithPricePrecisionString(d.StopLossPrice)
	takeProfitPrice := d.GetValueWithPricePrecisionString(d.TakeProfitPrice)

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
	msg = fmt.Sprintf("%sExpiration: %s", msg, expiration)
	return msg
}
