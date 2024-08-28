package exchange

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"gihub.com/shahinrahimi/teletradebot/models"
	"github.com/adshao/go-binance/v2/futures"
)

type BinanceClient struct {
	l          *log.Logger
	client     *futures.Client
	Symbols    []*futures.SymbolPrice
	ListenKey  string
	UseTestnet bool
}

func NewBinanceClient(l *log.Logger, apiKey string, secretKey string, useTestnet bool) *BinanceClient {
	futures.UseTestnet = useTestnet
	client := futures.NewClient(apiKey, secretKey)
	return &BinanceClient{l: l, client: client, UseTestnet: useTestnet}
}

func (b *BinanceClient) UpdateListenKey() error {
	b.l.Println(b.UseTestnet)
	listenKey, err := b.client.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		return err
	}
	b.ListenKey = listenKey
	return nil
}

func (b *BinanceClient) UserDataStream() {
	b.l.Println(b.UseTestnet)
	futures.UseTestnet = b.UseTestnet

	wsHandler := func(event *futures.WsUserDataEvent) {
		b.l.Println("we got a event", event)
		b.l.Println("we got a event", event.AccountConfigUpdate)
		b.l.Println("we got a event", event.AccountUpdate)
		b.l.Println("we got a event", event.CrossWalletBalance)
		b.l.Println("we got a event", event.Event)
		b.l.Println("we got a event", event.MarginCallPositions)
		b.l.Println("we got a event", event.OrderTradeUpdate)
		b.l.Println("we got a event", event.Time)
		b.l.Println("we got a event", event.TransactionTime)
		b.l.Println("we got a event", event.AccountUpdate.Balances)
		b.l.Println("we got a event", event.AccountUpdate.Positions)
		b.l.Println("we got a event", event.AccountUpdate.Reason)
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err := futures.WsUserDataServe(b.ListenKey, wsHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	b.l.Println("WebSocket connection established. Listening for events...")
	<-doneC

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
func (b *BinanceClient) GetSymbol(pair string) (*futures.SymbolPrice, error) {
	for _, symbol := range b.Symbols {
		if symbol.Symbol == pair {
			return symbol, nil
		}
	}
	return nil, fmt.Errorf("symbol not found")
}
func (b *BinanceClient) GetQuantity(o *models.Order) (string, error) {
	if err := b.UpdateTickers(); err != nil {
		return "", err
	}
	symbol, err := b.GetSymbol(o.Pair)
	if err != nil {
		return "", err
	}
	balanceStr, err := b.GetBalance()
	if err != nil {
		return "", err
	}
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return "", err
	}
	price, err := strconv.ParseFloat(symbol.Price, 64)
	if err != nil {
		return "", err
	}
	size := balance * float64(o.SizePercent) / 100
	fmt.Println(size, balance, o.SizePercent, float64(o.SizePercent)/100)
	return fmt.Sprintf("%.2f", size/price), nil
}
func (b *BinanceClient) GetBalance() (string, error) {
	balances, err := b.client.NewGetBalanceService().Do(context.Background())
	if err != nil {
		b.l.Printf("error getting account balance: %v", err)
		return "", err
	}
	for _, b := range balances {
		if b.Asset == "USDT" {
			return b.AvailableBalance, nil
		}
	}
	return "", fmt.Errorf("there is no balance found!")
}

func (b *BinanceClient) GetKline(o *models.Order) (*futures.Kline, error) {
	// TODO check if interval works for other intervals (1h and 4h works)
	klines, err := b.client.NewKlinesService().Limit(100).Interval(o.Candle).Symbol(o.Pair).Do(context.Background())
	if err != nil {
		return nil, err
	}
	// return before last element
	return klines[len(klines)-2], nil
}

func convertTime(timestamp int64) time.Time {
	// Convert milliseconds to seconds and nanoseconds
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1e6

	// Convert to time.Time
	return time.Unix(seconds, nanoseconds)
}

func (b *BinanceClient) PlaceOrder(o *models.Order) error {

	quantity, err := b.GetQuantity(o)
	if err != nil {
		return err
	}
	kline, err := b.GetKline(o)
	if err != nil {
		return err
	}
	h, err := strconv.ParseFloat(kline.High, 64)
	if err != nil {
		return err
	}
	l, err := strconv.ParseFloat(kline.Low, 64)
	if err != nil {
		return err
	}
	// for target
	// r := h-l
	//
	var side futures.SideType
	var stopPrice float64
	if o.Side == models.SIDE_L {
		side = futures.SideTypeBuy
		stopPrice = h + float64(o.Offset)
	} else {
		side = futures.SideTypeSell
		stopPrice = l - float64(o.Offset)
	}
	// should never happen
	if stopPrice <= 0 {
		return fmt.Errorf("price could not be zero or negative")
	}
	fmt.Printf("the quantity calculated for trade is %s", quantity)
	_, err = b.client.NewCreateOrderService().Symbol(o.Pair).Side(side).Type(futures.OrderTypeStopMarket).Quantity(quantity).StopPrice(fmt.Sprintf("%.2f", stopPrice)).Do(context.Background())
	if err != nil {
		// TODO add error handler if order type is not good for current price
		return err
	}
	// TODO add stoploss and take profit as well
	return nil
}
