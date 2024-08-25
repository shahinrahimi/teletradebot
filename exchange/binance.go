package exchange

import (
	"context"
	"log"

	"gihub.com/shahinrahimi/teletradebot/models"
	"github.com/adshao/go-binance/v2/futures"
)

type BinanceClient struct {
	l       *log.Logger
	client  *futures.Client
	Symbols []*futures.SymbolPrice
}

func NewBinanceClient(l *log.Logger, apiKey string, secretKey string) *BinanceClient {
	futures.UseTestnet = true
	client := futures.NewClient(apiKey, secretKey)

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

func (b *BinanceClient) PlaceOrder(o *models.Order) error {
	var side futures.SideType
	if o.Side == models.SIDE_L {
		side = futures.SideTypeBuy
	} else {
		side = futures.SideTypeSell
	}
	_, err := b.client.NewCreateOrderService().Symbol(o.Pair).Side(side).Type(futures.OrderTypeStopMarket).Quantity("1").StopPrice("1000.00").Do(context.Background())
	if err != nil {
		return err
	}
	return nil

}
