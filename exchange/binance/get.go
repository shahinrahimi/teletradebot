package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) getAvailableBalance() (float64, error) {
	if bc.lastSymbolPrices == nil {
		return 0, fmt.Errorf("latest account not available right now")
	}
	for _, balance := range bc.lastAccount.Assets {
		if balance.Asset == "USDT" {
			return strconv.ParseFloat(balance.AvailableBalance, 64)
		}
	}
	return 0, fmt.Errorf("no account information available")
}

func (bc *BinanceClient) CheckSymbol(symbol string) bool {
	if bc.lastExchangeInfo == nil {
		return false
	}
	for _, s := range bc.lastExchangeInfo.Symbols {
		if s.Symbol == symbol {
			return true
		}
	}
	return false
}

func (bc *BinanceClient) getSymbol(symbol string) (*futures.Symbol, error) {
	if bc.lastExchangeInfo == nil {
		return nil, fmt.Errorf("exchange info not available right now")
	}
	for _, s := range bc.lastExchangeInfo.Symbols {
		if s.Symbol == symbol {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("symbol %s not found", symbol)
}

func (bc *BinanceClient) getLatestPrice(symbol string) (float64, error) {
	if bc.lastSymbolPrices == nil {
		return 0, fmt.Errorf("latest prices not available right now")
	}
	for _, sp := range bc.lastSymbolPrices {
		if sp.Symbol == symbol {
			return strconv.ParseFloat(sp.Price, 64)
		}
	}
	return 0, fmt.Errorf("latest price not available for symbol %s", symbol)
}

func (bc *BinanceClient) getLastClosedKline(ctx context.Context, t *models.Trade) (*futures.Kline, error) {
	klines, err := bc.Client.NewMarkPriceKlinesService().
		Limit(100).
		Interval(t.Timeframe).
		Symbol(t.Symbol).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	// loop through klines and return the most recent completely closed candle
	for i := len(klines) - 1; i >= 0; i-- {
		candleCloseTime := utils.ConvertTime(klines[i].CloseTime)
		// check if close time in the past
		if (time.Until(candleCloseTime)) < 0 {
			return klines[i], nil
		}
	}

	return nil, fmt.Errorf("failed to locate before last candle")
}
