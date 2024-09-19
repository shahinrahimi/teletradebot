package bitmex

import (
	"context"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
)

func (mc *BitmexClient) PlaceTrade(ctx context.Context, t *models.Trade) (*swagger.Order, *PreparedOrder, error) {
	mc.l.Printf("executing order for trade ID: %d", t.ID)
	po, err := mc.prepareOrder(ctx, t)
	if err != nil {
		mc.l.Printf("trade could not be executed, error in preparing state: %v", err)
		return nil, nil, err
	}
	res, err := mc.PlaceOrder(ctx, po)
	return res, po, err
}
