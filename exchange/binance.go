package exchange

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
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
	listenKey, err := b.client.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		return err
	}
	b.ListenKey = listenKey
	return nil
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
func (b *BinanceClient) TrackOrder() error {
	orders, err := b.client.NewListOrdersService().Do(context.Background())
	if err != nil {
		return err
	}
	for _, order := range orders {
		b.l.Println(order.OrderID, order.Status, order.Symbol, order.StopPrice)
	}
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
func (b *BinanceClient) GetQuantity(t *models.Trade) (string, error) {
	if err := b.UpdateTickers(); err != nil {
		return "", err
	}
	symbol, err := b.GetSymbol(t.Pair)
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
	size := balance * float64(t.SizePercent) / 100
	quantity := size / price
	b.l.Printf("the balance is %.2f, the size of trade will be %.2f and the quantity will be %0.3f %s", balance, size, quantity, symbol.Symbol)
	return fmt.Sprintf("%.3f", quantity), nil
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
	return "", fmt.Errorf("there is no balance found")
}
func (b *BinanceClient) GetKline(t *models.Trade) (*futures.Kline, error) {
	// TODO check if interval works for other intervals (1h and 4h works)
	klines, err := b.client.NewKlinesService().Limit(100).Interval(t.Candle).Symbol(t.Pair).Do(context.Background())
	if err != nil {
		return nil, err
	}
	// return before last element
	return klines[len(klines)-2], nil
}
func (b *BinanceClient) PlaceOrder(t *models.Trade) (*futures.CreateOrderResponse, error) {

	quantity, err := b.GetQuantity(t)
	if err != nil {
		return nil, err
	}
	kline, err := b.GetKline(t)
	if err != nil {
		return nil, err
	}
	h, err := strconv.ParseFloat(kline.High, 64)
	if err != nil {
		return nil, err
	}
	l, err := strconv.ParseFloat(kline.Low, 64)
	if err != nil {
		return nil, err
	}
	// for target
	// r := h-l
	//
	var side futures.SideType
	var stopPrice float64
	if t.Side == types.SIDE_L {
		side = futures.SideTypeBuy
		stopPrice = h + float64(t.Offset)
	} else {
		side = futures.SideTypeSell
		stopPrice = l - float64(t.Offset)
	}
	// should never happen
	if stopPrice <= 0 {
		return nil, fmt.Errorf("price could not be zero or negative")
	}
	b.l.Printf("the quantity calculated for trade is %s", quantity)
	var order *futures.CreateOrderService
	order = b.client.NewCreateOrderService()
	order = order.Symbol(t.Pair).Side(side).Quantity(quantity)
	order = order.StopPrice(fmt.Sprintf("%.2f", stopPrice))
	order = order.Type(futures.OrderTypeStopMarket)
	order = order.WorkingType(futures.WorkingTypeMarkPrice)
	res, err := order.Do(context.Background())
	//res.OrderID
	if err != nil {
		// TODO add error handler if order type is not good for current price
		return nil, err
	}
	// TODO implement a function to cancel the order after desired duration
	return res, nil
}

func (b *BinanceClient) CancelOrder(orderID int64, symbol string) error {
	var order *futures.CancelOrderService
	order = b.client.NewCancelOrderService()
	order = order.OrderID(orderID).Symbol(symbol)
	_, err := order.Do(context.Background())
	if err != nil {
		b.l.Printf("error in canceling the order: %v", err)
		return err
	}
	return nil

}

func (b *BinanceClient) GetOrder(orderID int64, symbol string) (*futures.Order, error) {
	var order *futures.GetOrderService
	order = b.client.NewGetOrderService()
	order = order.OrderID(orderID).Symbol(symbol)
	res, err := order.Do(context.Background())
	if err != nil {
		b.l.Printf("error in getting the order: %v", err)
		return nil, err
	}
	return res, nil
}
