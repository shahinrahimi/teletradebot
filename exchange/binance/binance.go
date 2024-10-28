package binance

import (
	"context"
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

func (bc *BinanceClient) CheckSymbol(ctx context.Context, symbol string) (bool, error) {
	if _, err := bc.GetSymbol(ctx, symbol); err != nil {
		return false, err
	}
	return true, nil
	// fallback
}

func (bc *BinanceClient) CheckMultiAssetMode(ctx context.Context) (bool, error) {
	accountInfo, err := bc.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return false, fmt.Errorf("error getting account info: %v", err)
	}

	return accountInfo.MultiAssetsMargin, nil
}

func (bc *BinanceClient) GetSymbol(ctx context.Context, symbol string) (*futures.Symbol, error) {
	// fallback
	if bc.lastExchangeInfo == nil {
		ef, err := bc.client.NewExchangeInfoService().Do(ctx)
		if err != nil {
			bc.l.Printf("error fetching exchange info: %v", err)
			return nil, fmt.Errorf("error fetching exchange info: %v", err)
		}
		for _, s := range ef.Symbols {
			if s.Symbol == symbol {
				bc.lastExchangeInfo = ef
				return &s, nil
			}
		}
		return nil, fmt.Errorf("symbol %s not found for binance", symbol)
	} else {
		for _, s := range bc.lastExchangeInfo.Symbols {
			if s.Symbol == symbol {
				return &s, nil
			}
		}
		return nil, fmt.Errorf("symbol %s not found", symbol)
	}

}
