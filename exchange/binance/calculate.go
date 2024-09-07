package binance

import (
	"fmt"
	"math"
	"strconv"

	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
	"github.com/adshao/go-binance/v2/futures"
)

func (bc *BinanceClient) calculateStopPrice(t *models.Trade, k *futures.Kline, s *futures.Symbol) (string, error) {
	h, err := strconv.ParseFloat(k.High, 64)
	if err != nil {
		return "", err
	}
	l, err := strconv.ParseFloat(k.Low, 64)
	if err != nil {
		return "", err
	}
	var stopPrice float64

	if t.Side == types.SIDE_L {
		stopPrice = h + t.Offset
	} else {
		stopPrice = l - t.Offset
	}

	// Ensure stop price is positive
	if stopPrice <= 0 {
		return "", fmt.Errorf("price cannot be zero or negative")
	}

	pricePrecision := math.Pow10(int(-s.PricePrecision))
	stopPrice = math.Floor(stopPrice/pricePrecision) * pricePrecision

	return fmt.Sprintf("%.*f", s.PricePrecision, stopPrice), nil
}

func (bc *BinanceClient) calculateStopLossPrice(t *models.Trade, k *futures.Kline, s *futures.Symbol, basePrice string) (string, error) {
	h, err := strconv.ParseFloat(k.High, 64)
	if err != nil {
		return "", err
	}
	l, err := strconv.ParseFloat(k.Low, 64)
	if err != nil {
		return "", err
	}
	r := h - l
	sp, err := strconv.ParseFloat(basePrice, 64)
	if err != nil {
		return "", err
	}
	var sl float64
	if t.Side == types.SIDE_L {
		sl = sp - (r * (float64(t.StopLoss)) / 100)
	} else {
		sl = sp + (r * (float64(t.TakeProfit)) / 100)
	}

	pricePrecision := math.Pow10(int(-s.PricePrecision))
	sl = math.Floor(sl/pricePrecision) * pricePrecision

	stopLossPrice := fmt.Sprintf("%.*f", s.PricePrecision, sl)

	return stopLossPrice, nil
}

func (bc *BinanceClient) calculateTakeProfitPrice(t *models.Trade, k *futures.Kline, s *futures.Symbol, basePrice string) (string, error) {
	h, err := strconv.ParseFloat(k.High, 64)
	if err != nil {
		return "", err
	}
	l, err := strconv.ParseFloat(k.Low, 64)
	if err != nil {
		return "", err
	}
	r := h - l
	sp, err := strconv.ParseFloat(basePrice, 64)
	if err != nil {
		return "", err
	}
	var tp float64
	if t.Side == types.SIDE_L {
		tp = sp + (r * (float64(t.StopLoss)) / 100)
	} else {
		tp = sp - (r * (float64(t.TakeProfit)) / 100)
	}

	pricePrecision := math.Pow10(int(-s.PricePrecision))
	tp = math.Floor(tp/pricePrecision) * pricePrecision

	takeProfitPrice := fmt.Sprintf("%.*f", s.PricePrecision, tp)

	return takeProfitPrice, nil
}
