package bitmex

import (
	"context"
	"fmt"
	"time"

	"github.com/antihax/optional"
	"github.com/shahinrahimi/teletradebot/models"
	swagger "github.com/shahinrahimi/teletradebot/swagger"
)

var (
	bitmexTypes = map[string]string{
		"FFWCSX": "Perpetual Contracts",
		"FFWCSF": "Perpetual Futures",
		"IFXXXP": "Spot",
		"FFCCSX": "Futures",
		"MRBXXX": "BitMEX Basket Index",
		"MRCXXX": "BitMEX Crypto Index",
		"MRFXXX": "BitMEX FX Index",
		"MRRXXX": "BitMEX Lending/Premium Index",
		"MRIXXX": "BitMEX Volatility Index",
	}
)

func (mc *BitmexClient) CheckSymbol(symbol string) bool {
	// opts := &swagger.InstrumentApiInstrumentGetOpts{
	// 	Symbol: optional.NewString(symbol),
	// }
	instruments, _, err := mc.client.InstrumentApi.InstrumentGetActive(mc.auth)
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

func (mc *BitmexClient) GetMargins(ctx context.Context) ([]swagger.Margin, error) {
	ctx = mc.GetAuthContext(ctx)
	opts := swagger.UserApiUserGetMarginOpts{
		Currency: optional.NewString("all"),
	}
	margins, _, err := mc.client.UserApi.UserGetMargins(ctx, &opts)
	if err != nil {
		mc.l.Printf("failed to retrieve margins: %v", err)
		return nil, err
	}
	return margins, nil
}

func (mc *BitmexClient) GetBalance(ctx context.Context, currency string) (float64, error) {
	margins, err := mc.GetMargins(ctx)
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

func (mc *BitmexClient) GetBalanceXBt(ctx context.Context) (float64, error) {
	currency := "XBt"
	return mc.GetBalance(ctx, currency)
}

func (mc *BitmexClient) GetBalanceUSDt(ctx context.Context) (float64, error) {
	currency := "USDt"
	return mc.GetBalance(ctx, currency)
}

func (mc *BitmexClient) GetInstrument(ctx context.Context, t *models.Trade) (*swagger.Instrument, error) {
	ctx = mc.GetAuthContext(ctx)
	instruments, _, err := mc.client.InstrumentApi.InstrumentGetActive(ctx)
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

func (mc *BitmexClient) GetMarketPrice(ctx context.Context, t *models.Trade) (float64, error) {
	ctx = mc.GetAuthContext(ctx)
	opts := &swagger.InstrumentApiInstrumentGetOpts{
		Symbol: optional.NewString(t.Symbol),
	}
	instruments, _, err := mc.client.InstrumentApi.InstrumentGet(ctx, opts)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve market price: %v", err)
	}
	if len(instruments) > 0 {
		mc.l.Printf("lotsize: %0.1f  maxorderQty: %0.1f", instruments[0].LotSize, instruments[0].MaxOrderQty)
		return instruments[0].MarkPrice, nil
	}
	return 0, fmt.Errorf("no market data available for symbol: %s", t.Symbol)
}

func (mc *BitmexClient) GetLastClosedCandleOld(ctx context.Context, t *models.Trade) (*swagger.TradeBin, error) {
	ctx = mc.GetAuthContext(ctx)
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
		Filter:  optional.NewString(fitler),
	}
	// try until we get tradebucket with positive remaining time
	for i := 0; i < 10; i++ {
		tradeBins, _, err := mc.client.TradeApi.TradeGetBucketed(ctx, &params)
		if err != nil {
			mc.l.Printf("error: %v", err)
			mc.l.Printf("err type: %T", err)
			if ApiErr, ok := (err).(swagger.GenericSwaggerError); ok {
				mc.l.Printf("error creating a order: %s , errof: %s", string(ApiErr.Body()), ApiErr.Error())
			}
			return nil, err
		}
		for _, bin := range tradeBins {
			expiration, err := mc.calculateExpiration(t, bin.Timestamp)
			if err != nil && expiration == 0 {
				// lets try with delay sleep
				mc.l.Printf("retrying getting last closed candle with delay sleep 5s")
				time.Sleep(time.Second * 5)
				break
			} else if err != nil && expiration == -1 {
				return nil, err
			} else {
				return &bin, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to get last closed candle")

}
