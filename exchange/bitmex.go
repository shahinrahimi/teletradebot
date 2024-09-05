package exchange

import (
	"context"
	"fmt"
	"log"

	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/utils"
	"github.com/antihax/optional"
	"github.com/qct/bitmex-go/swagger"
)

type BitmexClient struct {
	l      *log.Logger
	client *swagger.APIClient
	ctx    context.Context
}

func NewBitmexClient(l *log.Logger, url string, apiKey string, apiSec string) *BitmexClient {
	cfg := swagger.NewConfiguration()
	cfg.BasePath = "https://testnet.bitmex.com/api/v1"
	client := swagger.NewAPIClient(cfg)
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    apiKey,
		Secret: apiSec,
	})

	return &BitmexClient{
		l:      l,
		client: client,
		ctx:    auth,
	}
}

func (mc *BitmexClient) GetKlineData() error {
	params := swagger.TradeApiTradeGetBucketedOpts{
		BinSize: optional.NewString("1h"),
		Count:   optional.NewFloat32(10),
		Reverse: optional.NewBool(true),
		Symbol:  optional.NewString("XBTUSD"),
		// StartTime: optional.NewTime(time.Now().Add(-24 * time.Hour)),
	}
	klineData, _, err := mc.client.TradeApi.TradeGetBucketed(mc.ctx, &params)
	if err != nil {
		utils.PrintStructFields(err)
		mc.l.Printf("err type: %T", err)
		if ApiErr, ok := (err).(swagger.GenericSwaggerError); ok {
			mc.l.Printf("error creating a order: %s , errof: %s", string(ApiErr.Body()), ApiErr.Error())
		}

		return err
	}
	// Print the Kline data
	for _, kline := range klineData {
		mc.l.Printf("Kline: %+v\n", kline)
	}
	// utils.PrintStructFields(klineData)
	return nil
}

func (mc *BitmexClient) GetWallet() error {
	params := &swagger.UserApiUserGetWalletOpts{
		Currency: optional.NewString("usdt"),
	}
	sw, resp, err := mc.client.UserApi.UserGetWallet(mc.ctx, params)
	if err != nil {
		utils.PrintStructFields(err)
		mc.l.Printf("err type: %T", err)
		if ApiErr, ok := (err).(swagger.GenericSwaggerError); ok {
			mc.l.Printf("error creating a order: %s , errof: %s", string(ApiErr.Body()), ApiErr.Error())
		}

		return err
	}
	utils.PrintStructFields(sw)
	utils.PrintStructFields(resp)
	return nil
}

func (mc *BitmexClient) GetAvailableBalance() (float64, error) {
	margin, _, err := mc.client.UserApi.UserGetMargin(mc.ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve margin: %v", err)
	}
	availableBalance := float64(margin.AvailableMargin) / 100000000 // Convert satoshis to BTC
	fmt.Println(availableBalance)
	return availableBalance, nil
}

func (mc *BitmexClient) GetMarketPrice() (float64, error) {

	instruments, _, err := mc.client.InstrumentApi.InstrumentGet(mc.ctx, &swagger.InstrumentApiInstrumentGetOpts{
		Symbol: optional.NewString("XBTUSD"),
		Count:  optional.NewFloat32(1),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve market price: %v", err)
	}
	if len(instruments) > 0 {
		mc.l.Println(instruments[0])
		return instruments[0].MarkPrice, nil
	}
	return 0, fmt.Errorf("no market data available for symbol: %s", "XBTUSD")
}

func (mc *BitmexClient) GetAvailableSymbols() {
	instruments, _, err := mc.client.InstrumentApi.InstrumentGet(nil, nil)
	if err != nil {
		fmt.Errorf("failed to retrieve instruments: %v", err)
	}

	// Extract symbols from the retrieved instruments
	var symbols []string
	for _, instrument := range instruments {
		symbols = append(symbols, instrument.Symbol)
		mc.l.Println(instrument.Symbol, instrument.MarkPrice, instrument.MaxOrderQty, instrument.MaxPrice)
	}
}

func (mc *BitmexClient) TryPlaceOrderForTrade(t *models.Trade) error {
	// orderApi := restful.NewOrderApi(mc.client.OrderApi, mc.ctx)
	//mc.GetWallet()
	mc.GetAvailableBalance()

	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString("Buy"),
		OrderQty: optional.NewFloat32(100),
		OrdType:  optional.NewString("Stop"),
		StopPx:   optional.NewFloat64(300000),
	}
	so, _, err := mc.client.OrderApi.OrderNew(mc.ctx, t.Symbol, params)
	// resp, orderId, err := orderApi.LimitBuy(t.Pair, 0.1, 30000.0, "4145153523413131314")
	if err != nil {
		utils.PrintStructFields(err)
		mc.l.Printf("err type: %T", err)
		if ApiErr, ok := (err).(swagger.GenericSwaggerError); ok {
			mc.l.Printf("error creating a order: %s , errof: %s", string(ApiErr.Body()), ApiErr.Error())
		}

		return err
	}
	mc.l.Printf("successfully creating a orderId: %s, orderId: %s", so.ClOrdID, so.OrderID)
	// mc.client.OrderApi.OrderNew(mc.ctx, t.Pair)
	// op := order.NewOrderNewParams().
	// 	WithSymbol(t.Pair).
	// 	WI
	// 	WithSimpleOrderQty(&quantity)

	// o, err := mc.client.Order.OrderNew(op)
	// if err != nil {
	// 	mc.l.Printf("err type: %T", err)
	// 	mc.l.Printf("error creating a order: %v", err)
	// 	if apiErr, ok := err.(*order.OrderNewBadRequest); ok {
	// 		mc.l.Printf("error %s", apiErr.Payload.Error.Message)

	// 	}
	// 	utils.PrintStructFields(err)
	// 	mc.l.Printf("err type: %T", err)
	// 	mc.l.Printf("error creating a order: %s", err)
	// 	return err
	// }
	// mc.l.Printf("successfully creating a orderId: %s", o.Payload.ClOrdID)
	return nil
}
