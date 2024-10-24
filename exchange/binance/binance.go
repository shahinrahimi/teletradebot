package binance

import (
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2/futures"
)

type BinanceClient struct {
	l                *log.Logger
	client           *futures.Client
	ListenKey        string
	UseTestnet       bool
	lastExchangeInfo *futures.ExchangeInfo
	ReverseEnabled   bool
	DbgChan          chan string
}

func NewBinanceClient(l *log.Logger, apiKey string, secretKey string, useTestnet bool, dbgChan chan string) *BinanceClient {
	futures.UseTestnet = useTestnet
	client := futures.NewClient(apiKey, secretKey)
	return &BinanceClient{
		l:              l,
		client:         client,
		UseTestnet:     useTestnet,
		ReverseEnabled: true,
		DbgChan:        dbgChan,
	}
}

func (bc *BinanceClient) CheckSymbol(symbol string) bool {
	if bc.lastExchangeInfo == nil {
		bc.l.Printf("exchange info not available right now please try after some time")
		return false
	}
	for _, s := range bc.lastExchangeInfo.Symbols {
		if s.Symbol == symbol {
			return true
		}
	}
	return false
}

func (bc *BinanceClient) GetSymbol(symbol string) (*futures.Symbol, error) {
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
