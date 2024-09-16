package bitmex

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/antihax/optional"
	"github.com/shahinrahimi/teletradebot/models"
	swagger "github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
)

type PreparedOrder struct {
	Symbol     string
	Quantity   float32
	StopPrice  float64
	Side       string
	Expiration time.Duration
}

func (mc *BitmexClient) PrepareOrder(ctx context.Context, t *models.Trade) (*PreparedOrder, error) {
	var p PreparedOrder
	p.Symbol = t.Symbol

	candle, err := mc.GetLastClosedCandle(t)
	if err != nil {
		mc.l.Printf("Error fetching last closed candle: %v", err)
		return nil, err
	}

	if t.Side == types.SIDE_L {
		p.Side = "Buy"
		p.StopPrice = candle.High + t.Offset
	} else {
		p.Side = "Sell"
		p.StopPrice = candle.Low - t.Offset
	}

	balance, err := mc.GetBalanceUSDt()
	if err != nil {
		mc.l.Printf("Error fetching balance: %v", err)
		return nil, err
	}
	mc.l.Printf("Fetched balance: %f", balance)

	instrument, err := mc.GetInstrument(t)
	if err != nil {
		mc.l.Printf("Error fetching instrument: %v", err)
		return nil, err
	}

	lotSize := instrument.LotSize
	quantity := (balance * (float64(t.Size) / 100000.0)) / instrument.MarkPrice
	if quantity < float64(lotSize) {
		return nil, fmt.Errorf("the calculated quantity (%.2f) less than the lotsize (%.1f)", quantity, lotSize)
	}
	roundedQuantity := math.Floor(quantity/float64(lotSize)) * float64(lotSize)
	mc.l.Printf("Calculated order quantity: %0.2f rounded (lotsize[%.1f]): %0.2f", quantity, lotSize, roundedQuantity)

	// if quantity < float64(instrument.LotSize) {
	// 	return nil, fmt.Errorf("the calculated quantity is less than instrument lot size")
	// }

	// candleDuration, err := types.GetDuration(t.Timeframe)
	// if err != nil {
	// 	return nil, err
	// }

	// candleCloseTime := utils.ConvertTime(candle.Timestamp)
	// remainingTime := candleDuration + time.Until(candleCloseTime)
	// if remainingTime < 0 {
	// 	return nil, fmt.Errorf("remaining time should not be negative number: %d", remainingTime)
	// }

	p.Quantity = float32(roundedQuantity)
	//p.Expiration = remainingTime

	return &p, nil

}

func (mc *BitmexClient) PlacePreparedOrder(p *PreparedOrder) (*swagger.Order, error) {

	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(p.Side),
		OrderQty: optional.NewFloat32(p.Quantity),
		OrdType:  optional.NewString("Stop"),
		StopPx:   optional.NewFloat64(p.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(mc.auth, p.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) PlaceOrder(ctx context.Context, p *PreparedOrder) (*swagger.Order, error) {
	ctx = context.WithValue(ctx, swagger.ContextAPIKey, swagger.APIKey{
		Key:    mc.ApiKey,
		Secret: mc.ApiSec,
	})
	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(p.Side),
		OrderQty: optional.NewFloat32(p.Quantity),
		OrdType:  optional.NewString("Stop"),
		StopPx:   optional.NewFloat64(p.StopPrice),
	}
	order, _, err := mc.client.OrderApi.OrderNew(ctx, p.Symbol, params)
	return &order, err
}

func (mc *BitmexClient) GetOrder(ctx context.Context, symbol string, orderID string) (*swagger.Order, error) {

	params := &swagger.OrderApiOrderGetOrdersOpts{
		Symbol: optional.NewString(symbol),
	}
	orders, _, err := mc.client.OrderApi.OrderGetOrders(ctx, params)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		if o.OrderID == orderID {
			return &o, nil
		}
	}
	return nil, fmt.Errorf("order not found")
}
