package bitmex

import (
	"fmt"
	"math"
	"time"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (mc *BitmexClient) calculateQuantity(t *models.Trade, balance, markPrice, lotSize float64) (float64, error) {
	var quantity float64
	if markPrice == 0 {
		return 0, fmt.Errorf("mark price cannot be zero")
	}
	if lotSize == 0 {
		return 0, fmt.Errorf("lot size cannot be zero")
	}
	quantity = (balance * (float64(t.Size) / 100000.0)) / markPrice
	quantity = math.Floor(quantity/lotSize) * lotSize
	if quantity < lotSize {
		return 0, fmt.Errorf("the calculated quantity (%.2f) less than the lotsize (%.1f)", quantity, lotSize)
	}
	if mc.Verbose {
		mc.l.Printf("calculated quantity: %f", quantity)
	}
	return quantity, nil
}

func (mc *BitmexClient) calculateExpiration(t *models.Trade, closeTimestamp time.Time) (time.Duration, error) {
	var expiration time.Duration
	candleDuration, err := types.GetDuration(t.Timeframe)
	if err != nil {
		return -1, err
	}
	expiration = candleDuration + time.Until(closeTimestamp)
	if expiration < 0 {
		return 0, fmt.Errorf("expiration time should not be negative number: %d", expiration)
	}
	if mc.Verbose {
		mc.l.Printf("calculated expiration: %s", utils.FriendlyDuration(expiration))
	}
	return expiration, nil
}
func (mc *BitmexClient) calculateStopPrice(t *models.Trade, high, low, tickSize float64) (float64, error) {
	var stopPrice float64
	if tickSize == 0 {
		return 0, fmt.Errorf("tick size cannot be zero")
	}
	if t.Side == types.SIDE_L {
		stopPrice = high + t.Offset + 0.005*high
	} else {
		stopPrice = low - t.Offset - 0.005*low
	}

	ticks := stopPrice / tickSize
	stopPrice = math.Round(ticks) * tickSize
	if stopPrice <= 0 {
		return 0, fmt.Errorf("price cannot be zero or negative")
	}
	if mc.Verbose {
		mc.l.Printf("calculated stop price: %f", stopPrice)
	}
	return stopPrice, nil
}

func (mc *BitmexClient) calculateStopLossPrice(t *models.Trade, high, low, tickSize, basePrice float64) (float64, error) {
	var stopPrice float64
	r := high - low
	if t.Side == types.SIDE_L {
		stopPrice = basePrice - (r * (float64(t.StopLoss)) / 100)
	} else {
		stopPrice = basePrice + (r * (float64(t.StopLoss)) / 100)
	}
	ticks := stopPrice / tickSize
	stopPrice = math.Round(ticks) * tickSize
	if stopPrice <= 0 {
		return 0, fmt.Errorf("price cannot be zero or negative")
	}
	if mc.Verbose {
		mc.l.Printf("calculated stop loss price: %f", stopPrice)
	}
	return stopPrice, nil
}

func (mc *BitmexClient) calculateTakeProfitPrice(t *models.Trade, high, low, tickSize, basePrice float64) (float64, error) {
	var stopPrice float64
	r := high - low
	if t.Side == types.SIDE_L {
		stopPrice = basePrice + (r * (float64(t.TakeProfit)) / 100)
	} else {
		stopPrice = basePrice - (r * (float64(t.TakeProfit)) / 100)
	}
	ticks := stopPrice / tickSize
	stopPrice = math.Round(ticks) * tickSize
	if stopPrice <= 0 {
		return 0, fmt.Errorf("price cannot be zero or negative")
	}
	if mc.Verbose {
		mc.l.Printf("calculated take profit price: %f", stopPrice)
	}
	return stopPrice, nil
}
