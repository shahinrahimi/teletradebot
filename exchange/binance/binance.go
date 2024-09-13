package binance

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/store"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

type BinanceClient struct {
	l                *log.Logger
	s                store.Storage
	Client           *futures.Client
	ListenKey        string
	UseTestnet       bool
	lastSymbolPrices []*futures.SymbolPrice
	lastExchangeInfo *futures.ExchangeInfo
	lastAccount      *futures.Account
	MsgChan          chan types.BotMessage
}

func NewBinanceClient(l *log.Logger, s store.Storage, apiKey string, secretKey string, useTestnet bool, msgChan chan types.BotMessage) *BinanceClient {
	futures.UseTestnet = useTestnet
	client := futures.NewClient(apiKey, secretKey)
	return &BinanceClient{
		l:          l,
		s:          s,
		Client:     client,
		UseTestnet: useTestnet,
		MsgChan:    msgChan,
	}
}

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

func (bc *BinanceClient) getSymbol(t *models.Trade) (*futures.Symbol, error) {
	if bc.lastExchangeInfo == nil {
		return nil, fmt.Errorf("exchange info not available right now")
	}
	for _, s := range bc.lastExchangeInfo.Symbols {
		if s.Symbol == t.Symbol {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("symbol %s not found", t.Symbol)
}

func (bc *BinanceClient) getLatestPrice(t *models.Trade) (float64, error) {
	if bc.lastSymbolPrices == nil {
		return 0, fmt.Errorf("latest prices not available right now")
	}
	for _, sp := range bc.lastSymbolPrices {
		if sp.Symbol == t.Symbol {
			return strconv.ParseFloat(sp.Price, 64)
		}
	}
	return 0, fmt.Errorf("latest price not available for symbol %s", t.Symbol)
}

func (bc *BinanceClient) getQuantity(t *models.Trade) (string, error) {
	s, err := bc.getSymbol(t)
	if err != nil {
		return "", err
	}
	b, err := bc.getAvailableBalance()
	if err != nil {
		return "", err
	}
	p, err := bc.getLatestPrice(t)
	if err != nil {
		return "", err
	}

	size := b * float64(t.Size) / 100
	quantity := size / p

	// adjust quantity based on symbol precision
	quantityPrecision := math.Pow10(int(-s.QuantityPrecision))
	quantity = math.Floor(quantity/quantityPrecision) * quantityPrecision

	return fmt.Sprintf("%.*f", s.QuantityPrecision, quantity), nil
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

func (bc *BinanceClient) handleError(err error, userID int64) {
	if apiErr, ok := err.(*common.APIError); ok {
		msg := fmt.Sprintf("Binance API:\n\nMessage: %s\nCode: %d", apiErr.Message, apiErr.Code)
		bc.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
	} else {
		bc.l.Printf("error casting error to Api error type: %T", err)
	}
}
