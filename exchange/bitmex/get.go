package bitmex

import (
	"fmt"
	"time"

	swagger "gihub.com/shahinrahimi/teletradebot/go-client"
	"gihub.com/shahinrahimi/teletradebot/models"
	"github.com/antihax/optional"
)

func (mc *BitmexClient) CheckSymbol(symbol string) bool {
	opts := &swagger.InstrumentApiInstrumentGetOpts{
		Symbol: optional.NewString(symbol),
	}
	instruments, _, err := mc.client.InstrumentApi.InstrumentGet(mc.auth, opts)
	if err != nil {
		return false
	}
	for _, i := range instruments {
		if i.Symbol == symbol {
			return true
		}
	}
	return false
}

func (mc *BitmexClient) GetMargins() ([]swagger.Margin, error) {
	opts := swagger.UserApiUserGetMarginOpts{
		Currency: optional.NewString("all"),
	}
	margins, _, err := mc.client.UserApi.UserGetMargins(mc.auth, &opts)
	if err != nil {
		mc.l.Printf("failed to retrieve margins: %v", err)
		return nil, err
	}
	return margins, nil
}

func (mc *BitmexClient) GetBalance(currency string) (float64, error) {
	margins, err := mc.GetMargins()
	if err != nil {
		return 0, err
	}
	for _, m := range margins {
		if m.Currency == currency {
			return float64(m.AvailableMargin), nil
		}
	}
	return 0, fmt.Errorf("the currency '%s' not found", currency)
}

func (mc *BitmexClient) GetBalanceXBt() (float64, error) {
	currency := "XBt"
	return mc.GetBalance(currency)
}

func (mc *BitmexClient) GetBalanceUSDt() (float64, error) {
	currency := "USDt"
	return mc.GetBalance(currency)
}

func (mc *BitmexClient) GetInstrument(t *models.Trade) (*swagger.Instrument, error) {
	instruments, _, err := mc.client.InstrumentApi.InstrumentGetActive(mc.auth)
	if err != nil {
		return nil, err
	}
	for _, i := range instruments {
		if i.Symbol == t.Symbol {
			return &i, nil
		}
	}
	return nil, fmt.Errorf("could not find instrument")
}

func (mc *BitmexClient) GetMarketPrice(t *models.Trade) (float64, error) {
	opts := &swagger.InstrumentApiInstrumentGetOpts{
		Symbol: optional.NewString(t.Symbol),
	}
	instruments, _, err := mc.client.InstrumentApi.InstrumentGet(mc.auth, opts)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve market price: %v", err)
	}
	if len(instruments) > 0 {
		mc.l.Printf("lotsize: %0.1f  maxorderQty: %0.1f", instruments[0].LotSize, instruments[0].MaxOrderQty)
		return instruments[0].MarkPrice, nil
	}
	return 0, fmt.Errorf("no market data available for symbol: %s", t.Symbol)
}

func (mc *BitmexClient) GetLastClosedCandle(t *models.Trade) (*swagger.TradeBin, error) {
	endTime := time.Now().UTC().Format("2006-01-02 15:04")
	fitler := fmt.Sprintf(`{"endTime": "%s"}`, endTime)
	// TODO the available bin size 1m,5m,1h,1d
	// it should flexible and convert this binSiz to for example 1m => 3m or 5m or 5m to 15min
	params := swagger.TradeApiTradeGetBucketedOpts{
		BinSize: optional.NewString(t.Timeframe),
		Symbol:  optional.NewString(t.Symbol),
		Count:   optional.NewFloat32(1),
		// Partial: optional.NewBool(true),
		Reverse: optional.NewBool(true),
		// EndTime: optional.NewTime(time.Now().Add(-time.Hour * 100).UTC().Format(time.RFC3339)),
		Filter: optional.NewString(fitler),
		// Count:   optional.NewFloat32(1),
	}
	tradeBins, _, err := mc.client.TradeApi.TradeGetBucketed(mc.auth, &params)
	if err != nil {
		mc.l.Printf("error: %v", err)
		mc.l.Printf("err type: %T", err)
		if ApiErr, ok := (err).(swagger.GenericSwaggerError); ok {
			mc.l.Printf("error creating a order: %s , errof: %s", string(ApiErr.Body()), ApiErr.Error())
		}
		return nil, err
	}
	for index, bin := range tradeBins {
		mc.l.Printf("index: %d symbol: %s open: %0.2f close: %0.2f", index, bin.Symbol, bin.Open, bin.Close, bin.Timestamp.UTC())
	}
	if len(tradeBins) > 0 {
		return &tradeBins[0], nil
	}
	return nil, nil
}
