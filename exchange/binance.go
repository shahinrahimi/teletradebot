package exchange

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
	"gihub.com/shahinrahimi/teletradebot/utils"
	"github.com/adshao/go-binance/v2/futures"
)

type BinanceClient struct {
	l                  *log.Logger
	client             *futures.Client
	SymbolPriceChan    chan []*futures.SymbolPrice
	ExchangeInfoChan   chan *futures.ExchangeInfo
	AccountChan        chan *futures.Account
	accountTriggerChan chan struct{} // Channel to trigger immediate balance refresh
	ListenKey          string
	UseTestnet         bool

	// Fields to store the latest values
	LastSymbolPrices []*futures.SymbolPrice
	LastExchangeInfo *futures.ExchangeInfo
	LastAccount      *futures.Account
}

func NewBinanceClient(l *log.Logger, apiKey string, secretKey string, useTestnet bool) *BinanceClient {
	futures.UseTestnet = useTestnet
	client := futures.NewClient(apiKey, secretKey)
	return &BinanceClient{
		l:                  l,
		client:             client,
		SymbolPriceChan:    make(chan []*futures.SymbolPrice),
		ExchangeInfoChan:   make(chan *futures.ExchangeInfo),
		AccountChan:        make(chan *futures.Account),
		accountTriggerChan: make(chan struct{}, 1),
		UseTestnet:         useTestnet,
	}
}

func (bc *BinanceClient) GetLatestSymbolPrices() []*futures.SymbolPrice {
	return bc.LastSymbolPrices
}

func (bc *BinanceClient) GetLatestExchangeInfo() *futures.ExchangeInfo {
	return bc.LastExchangeInfo
}

func (bc *BinanceClient) GetLatestAccount() *futures.Account {
	return bc.LastAccount
}

func (bc *BinanceClient) StartPolling() {
	// Weight: 5
	go bc.PollAccount(time.Minute)
	// Weight: 1
	go bc.PollExchangeInfo(time.Minute)
	// Weight: 2
	go bc.PollSymbolPrice(time.Minute)
}

func (bc *BinanceClient) PollSymbolPrice(interval time.Duration) {
	for {
		res, err := bc.client.NewListPricesService().Do(context.Background())
		if err != nil {
			bc.l.Printf("Error fetching prices: %v", err)
			continue
		}
		bc.LastSymbolPrices = res
		//bc.SymbolPriceChan <- res
		time.Sleep(interval)
	}
}

func (bc *BinanceClient) PollExchangeInfo(interval time.Duration) {
	for {
		res, err := bc.client.NewExchangeInfoService().Do(context.Background())
		if err != nil {
			bc.l.Printf("Error fetching exchange info: %v", err)
			continue
		}
		bc.LastExchangeInfo = res
		//bc.l.Printf("server time %s", utils.FormatTimestamp(bc.LastExchangeInfo.ServerTime))
		time.Sleep(interval)
	}
}

func (bc *BinanceClient) PollAccount(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Immediate run before entering the loop
	bc.fetchAndSendAccount()

	for {
		select {
		case <-ticker.C: // regular interval trigger
			bc.fetchAndSendAccount()
		case <-bc.accountTriggerChan: // forced refresh trigger
			bc.fetchAndSendAccount()
		}
	}
}

func (bc *BinanceClient) fetchAndSendAccount() {
	res, err := bc.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		bc.l.Printf("Error fetching account: %v", err)
		return
	}
	bc.LastAccount = res
	//bc.AccountChan <- res
}

func (bc *BinanceClient) ForceAccountRefresh() {
	select {
	case bc.accountTriggerChan <- struct{}{}:
		// trigger sent
	default:
		// If there is already a trigger pending, do nothing
	}
}

func (bc *BinanceClient) GetLatestPrice(symbol string) (float64, error) {
	sps := bc.GetLatestSymbolPrices()
	for _, sp := range sps {
		if sp.Symbol == symbol {
			return strconv.ParseFloat(sp.Price, 64)
		}
	}
	return 0, fmt.Errorf("latest price not available for symbol %s", symbol)
}
func (bc *BinanceClient) GetSymbol(symbol string) (*futures.Symbol, error) {
	ei := bc.GetLatestExchangeInfo()
	for _, s := range ei.Symbols {
		if s.Symbol == symbol {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("symbol %s not found", symbol)
}
func (bc *BinanceClient) GetAvailableBalance() (float64, error) {
	account := bc.GetLatestAccount()
	if account == nil {
		return 0, fmt.Errorf("acount is nil")
	}
	for _, balance := range account.Assets {
		if balance.Asset == "USDT" {
			return strconv.ParseFloat(balance.AvailableBalance, 64)
		}
	}
	return 0, fmt.Errorf("no account information available")
}
func (bc *BinanceClient) GetQuantity(t *models.Trade) (string, error) {
	symbol, err := bc.GetSymbol(t.Pair)
	if err != nil {
		return "", err
	}
	balance, err := bc.GetAvailableBalance()
	if err != nil {
		return "", err
	}
	price, err := bc.GetLatestPrice(t.Pair)
	if err != nil {
		return "", err
	}

	size := balance * float64(t.SizePercent) / 100
	quantity := size / price

	// adjust quantity based on symbol precision
	quantityPrecision := math.Pow10(int(-symbol.QuantityPrecision))
	quantity = math.Floor(quantity/quantityPrecision) * quantityPrecision

	return fmt.Sprintf("%.*f", symbol.QuantityPrecision, quantity), nil
}
func (bc *BinanceClient) GetKlineBeforeLast(symbol string, candle string) (*futures.Kline, error) {
	klines, err := bc.client.NewMarkPriceKlinesService().
		Limit(100).
		Interval(candle).
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return klines[len(klines)-2], nil
}
func (bc *BinanceClient) GetStopPrice(t *models.Trade, kline *futures.Kline) (string, error) {
	h, err := strconv.ParseFloat(kline.High, 64)
	if err != nil {
		return "", err
	}
	l, err := strconv.ParseFloat(kline.Low, 64)
	if err != nil {
		return "", err
	}
	var stopPrice float64

	if t.Side == types.SIDE_L {
		stopPrice = h + t.Offset
	} else {
		stopPrice = l - t.Offset
	}
	fmt.Println(h, l, stopPrice)

	// Ensure stop price is positive
	if stopPrice <= 0 {
		return "", fmt.Errorf("price cannot be zero or negative")
	}

	// Adjust the stop price according to price precision
	symbol, err := bc.GetSymbol(t.Pair)
	if err != nil {
		return "", err
	}

	pricePrecision := math.Pow10(int(-symbol.PricePrecision))
	stopPrice = math.Floor(stopPrice/pricePrecision) * pricePrecision

	return fmt.Sprintf("%.*f", symbol.PricePrecision, stopPrice), nil
}

func (bc *BinanceClient) GetStopLossTakeProfitPrice(t *models.Trade, stopPrice string) (stopLossPrice string, takeProfitPrice string, err error) {
	kline, err := bc.GetKlineBeforeLast(t.Pair, t.Candle)
	if err != nil {
		return "", "", err
	}
	h, err := strconv.ParseFloat(kline.High, 64)
	if err != nil {
		return "", "", err
	}
	l, err := strconv.ParseFloat(kline.Low, 64)
	if err != nil {
		return "", "", err
	}
	r := h - l
	sp, err := strconv.ParseFloat(stopPrice, 64)
	if err != nil {
		return "", "", err
	}

	var sl float64
	var tp float64
	if t.Side == types.SIDE_L {
		sl = sp - (r * (float64(t.SLPercent)) / 100)
		tp = sp + (r * (float64(t.TPPercent)) / 100)
	} else {
		sl = sp + (r * (float64(t.SLPercent)) / 100)
		tp = sp - (r * (float64(t.TPPercent)) / 100)
	}

	// Adjust the stop price according to price precision
	symbol, err := bc.GetSymbol(t.Pair)
	if err != nil {
		return "", "", err
	}
	pricePrecision := math.Pow10(int(-symbol.PricePrecision))
	sl = math.Floor(sl/pricePrecision) * pricePrecision
	tp = math.Floor(tp/pricePrecision) * pricePrecision

	stopLossPrice = fmt.Sprintf("%.*f", symbol.PricePrecision, sl)
	takeProfitPrice = fmt.Sprintf("%.*f", symbol.PricePrecision, tp)

	return stopLossPrice, takeProfitPrice, nil
}

func (bc *BinanceClient) TryPlaceOrderForTrade(t *models.Trade) (*futures.CreateOrderResponse, error) {
	quantity, err := bc.GetQuantity(t)
	if err != nil {
		return nil, err
	}
	kline, err := bc.GetKlineBeforeLast(t.Pair, t.Candle)
	if err != nil {
		return nil, err
	}
	stopPrice, err := bc.GetStopPrice(t, kline)
	if err != nil {
		return nil, err
	}
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeBuy
	} else {
		side = futures.SideTypeSell
	}

	// Determine the remaining time until the candle closes
	candleDuration, err := types.GetDuration(t.Candle)
	if err != nil {
		return nil, err
	}

	candleCloseTime := kline.CloseTime
	remainingTime := candleDuration + time.Until(utils.ConvertTime(candleCloseTime))
	bc.l.Printf("candle duration: %s", utils.FriendlyDuration(candleDuration))
	bc.l.Printf("candle open time: %s", utils.ConvertTime(kline.OpenTime))
	bc.l.Printf("candle close time: %s", utils.ConvertTime(kline.CloseTime))
	bc.l.Printf("remaining duration: %s", utils.FriendlyDuration(remainingTime))
	bc.l.Printf("Placing %s order with quantity %s and stop price %s expires in: %s", side, quantity, stopPrice, utils.FriendlyDuration(remainingTime))
	if remainingTime < 0 {
		return nil, fmt.Errorf("remaining time should not be negative number: %d", remainingTime)
	}

	order := bc.client.NewCreateOrderService().
		Symbol(t.Pair).
		Side(side).
		Quantity(quantity).
		StopPrice(stopPrice).
		Type(futures.OrderTypeStopMarket).
		WorkingType(futures.WorkingTypeMarkPrice)

	res, err := order.Do(context.Background())
	if err != nil {
		return nil, err
	}

	// Start a timer to cancel the order if not filled in the remaining time
	go bc.scheduleOrderCancellation(res.OrderID, res.Symbol, remainingTime)

	return res, nil
}

func (bc *BinanceClient) TryPlaceStopLossAndTakeProfitTrade(t *models.Trade, orderUpdate *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, *futures.CreateOrderResponse, error, error) {
	stopLossPrice, takeProfitPrice, err := bc.GetStopLossTakeProfitPrice(t, orderUpdate.StopPrice)
	if err != nil {
		return nil, nil, err, err
	}
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}
	orgOrder := bc.client.NewCreateOrderService().Symbol(t.Pair).Side(side).Quantity(orderUpdate.OriginalQty).WorkingType(futures.WorkingTypeMarkPrice)
	slOrder := orgOrder.Type(futures.OrderTypeStopMarket).StopPrice(stopLossPrice)
	tpOrder := orgOrder.Type(futures.OrderTypeTakeProfitMarket).StopPrice(takeProfitPrice)
	// execute slOrder
	res1, err1 := slOrder.Do(context.Background())
	// execute tpOrder
	res2, err2 := tpOrder.Do(context.Background())
	if err1 != nil || err2 != nil {
		return nil, nil, err1, err2
	}
	return res1, res2, err1, err2
}

func (bc *BinanceClient) CancelOrder(orderID int64, symbol string) (*futures.CancelOrderResponse, error) {
	cancelOrder := bc.client.NewCancelOrderService().OrderID(orderID).Symbol(symbol)
	res, err := cancelOrder.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (bc *BinanceClient) scheduleOrderCancellation(orderID int64, symbol string, delay time.Duration) {
	time.AfterFunc(delay, func() {
		_, err := bc.CancelOrder(orderID, symbol)
		if err != nil {
			bc.l.Printf("Failed to cancel order %d: %v", orderID, err)
		} else {
			bc.l.Printf("order %d canceled after %v", orderID, delay)
		}
	})
}

func (b *BinanceClient) UpdateListenKey() error {
	listenKey, err := b.client.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		return err
	}
	b.ListenKey = listenKey
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
