package bitmex

import (
	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
	"github.com/antihax/optional"
	"github.com/qct/bitmex-go/swagger"
)

func (mc *BitmexClient) PlaceOrder(t *models.Trade) error {
	balance, err := mc.GetBalanceUSD()
	if err != nil {
		mc.l.Printf("Error fetching balance: %v", err)
		return err
	}
	mc.l.Printf("Fetched balance: %f", balance)

	price, err := mc.GetMarketPrice(t)
	if err != nil {
		mc.l.Printf("Error fetching market price: %v", err)
		return err
	}
	mc.l.Printf("Fetched market price: %f", price)

	candle, err := mc.GetLastClosedCandle(t)
	if err != nil {
		mc.l.Printf("Error fetching last closed candle: %v", err)
		return err
	}
	mc.l.Printf("Fetched last closed candle: %+v", candle)

	orderQ := (50000 * balance * (float64(t.SizePercent) / 100.0)) / price
	mc.l.Printf("Calculated order quantity: %f", orderQ)

	if orderQ < 1 {
		orderQ = 1
	}

	var side string
	var sp float64
	if t.Side == types.SIDE_L {
		side = "Buy"
		sp = candle.High + t.Offset
	} else {
		side = "Sell"
		sp = candle.Low - t.Offset
	}
	mc.l.Printf("Order side: %s, Stop price: %f", side, sp)

	params := &swagger.OrderApiOrderNewOpts{
		Side:     optional.NewString(side),
		OrderQty: optional.NewFloat32(float32(orderQ)),
		OrdType:  optional.NewString("Stop"),
		StopPx:   optional.NewFloat64(sp),
	}
	mc.l.Printf("Order parameters: %+v", params)

	order, _, err := mc.client.OrderApi.OrderNew(mc.auth, t.Symbol, params)
	if err != nil {
		mc.l.Printf("Error placing order: %v", err)
		return err
	}
	mc.l.Printf("Order placed successfully: %v", order.OrderID)
	return nil

}
