package binance

import (
	"log"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/types"
)

type BinanceClient struct {
	l                *log.Logger
	Client           *futures.Client
	ListenKey        string
	UseTestnet       bool
	lastSymbolPrices []*futures.SymbolPrice
	lastExchangeInfo *futures.ExchangeInfo
	lastAccount      *futures.Account
	MsgChan          chan types.BotMessage
}

func NewBinanceClient(l *log.Logger, apiKey string, secretKey string, useTestnet bool, msgChan chan types.BotMessage) *BinanceClient {
	futures.UseTestnet = useTestnet
	client := futures.NewClient(apiKey, secretKey)
	return &BinanceClient{
		l:          l,
		Client:     client,
		UseTestnet: useTestnet,
		MsgChan:    msgChan,
	}
}

// func (bc *BinanceClient) getQuantity(t *models.Trade) (string, error) {
// 	s, err := bc.getSymbol(t)
// 	if err != nil {
// 		return "", err
// 	}
// 	b, err := bc.getAvailableBalance()
// 	if err != nil {
// 		return "", err
// 	}
// 	p, err := bc.getLatestPrice(t)
// 	if err != nil {
// 		return "", err
// 	}

// 	size := b * float64(t.Size) / 100
// 	quantity := size / p

// 	// adjust quantity based on symbol precision
// 	quantityPrecision := math.Pow10(int(-s.QuantityPrecision))
// 	quantity = math.Floor(quantity/quantityPrecision) * quantityPrecision

// 	return fmt.Sprintf("%.*f", s.QuantityPrecision, quantity), nil
// }
