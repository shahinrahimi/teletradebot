package bitmexold

import (
	"context"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
)

func (mc *BitmexClient) PlaceTrade(ctx context.Context, t *models.Trade) (*swagger.Order, *models.Describer, error) {
	mc.l.Printf("executing order for trade ID: %d", t.ID)
	d, err := mc.FetchDescriber(ctx, t)
	if err != nil {
		mc.l.Printf("error fetching the describer %v", err)
		return nil, nil, err
	}
	po, err := mc.prepareMainOrder(ctx, d)
	if err != nil {
		mc.l.Printf("trade could not be executed, error in preparing state: %v", err)
		return nil, nil, err
	}
	res, err := mc.PlaceOrder(ctx, po)
	return res, d, err
}

func (mc *BitmexClient) PlaceTradeSLOrder(ctx context.Context, t *models.Trade, d *models.Describer, od *OrderData) (*swagger.Order, error) {
	mc.l.Printf("executing stop-loss order for trade ID: %d", t.ID)
	po := mc.prepareStopLossOrder(d, od)
	res, err := mc.PlaceSLOrder(ctx, po)
	return res, err
}

func (mc *BitmexClient) PlaceTradeTPOrder(ctx context.Context, t *models.Trade, d *models.Describer, od *OrderData) (*swagger.Order, error) {
	mc.l.Printf("executing take-profit order for trade ID: %d", t.ID)
	po := mc.prepareTakeProfitOrder(d, od)
	res, err := mc.PlaceTPOrder(ctx, po)
	return res, err
}

func (mc *BitmexClient) PlaceTradeReverseMainOrder(ctx context.Context, t *models.Trade, d *models.Describer, od *OrderData) (*swagger.Order, error) {
	mc.l.Printf("executing reverse main order for trade ID: %d", t.ID)
	po := mc.prepareReverseMainOrder(d, od)
	res, err := mc.PlaceOrder(ctx, po)
	return res, err
}

func (mc *BitmexClient) PlaceTradeReverseStopLossOrder(ctx context.Context, t *models.Trade, d *models.Describer, od *OrderData) (*swagger.Order, error) {
	mc.l.Printf("executing reverse stop-loss order for trade ID: %d", t.ID)
	po := mc.prepareReverseStopLossOrder(d, od)
	res, err := mc.PlaceSLOrder(ctx, po)
	return res, err
}

func (mc *BitmexClient) PlaceTradeReverseTakeProfitOrder(ctx context.Context, t *models.Trade, d *models.Describer, od *OrderData) (*swagger.Order, error) {
	mc.l.Printf("executing reverse take-profit order for trade ID: %d", t.ID)
	po := mc.prepareReverseTakeProfitOrder(d, od)
	res, err := mc.PlaceTPOrder(ctx, po)
	return res, err
}
