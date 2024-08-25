package exchange

import (
	"context"
	"log"

	"github.com/adshao/go-binance/v2"
)

type BinanceClient struct {
	l       *log.Logger
	client  *binance.Client
	Symbols []*binance.SymbolPrice
}

func NewBinanceClient(l *log.Logger, apiKey string, secretKey string) *BinanceClient {
	client := binance.NewClient(apiKey, secretKey)
	return &BinanceClient{l: l, client: client}
}
func (b *BinanceClient) UpdateTickers() error {
	symbols, err := b.client.NewListPricesService().Do(context.Background())
	if err != nil {
		b.l.Printf("error listing tickers: %v", err)
		return err
	}
	b.Symbols = symbols
	return nil
}
