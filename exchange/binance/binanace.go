package binance

import (
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/types"
)

type BinanceClient struct {
	l                *log.Logger
	client           *futures.Client
	ListenKey        string
	UseTestnet       bool
	lastExchangeInfo *futures.ExchangeInfo
	MsgChan          chan types.BotMessage
	ReverseEnabled   bool
}

func NewBinanceClient(l *log.Logger, apiKey string, secretKey string, useTestnet bool, msgChan chan types.BotMessage) *BinanceClient {
	futures.UseTestnet = useTestnet
	client := futures.NewClient(apiKey, secretKey)
	return &BinanceClient{
		l:              l,
		client:         client,
		UseTestnet:     useTestnet,
		MsgChan:        msgChan,
		ReverseEnabled: true,
	}
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
