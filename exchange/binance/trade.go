package binance

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
)

func (bc *BinanceClient) PlaceTrade(ctx context.Context, t *models.Trade) (*futures.CreateOrderResponse, *PreparedOrder, error) {
	bc.l.Printf("executing order for trade ID: %d", t.ID)
	po, err := bc.prepareOrder(ctx, t)
	if err != nil {
		bc.l.Printf("trade could not be executed, error in preparing state: %v", err)
		return nil, nil, err
	}
	res, err := bc.PlaceOrder(ctx, po)
	return res, po, err
}

func (bc *BinanceClient) PlaceTradeSLOrder(ctx context.Context, t *models.Trade, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, *PreparedOrder, error) {
	bc.l.Printf("executing stop-loss order for trade ID: %d", t.ID)
	po, err := bc.prepareSLOrder(ctx, t, f)
	if err != nil {
		bc.l.Printf("error during stop-loss order preparation: %v", err)
		return nil, nil, err
	}
	res, err := bc.placeSLOrder(ctx, po)
	return res, po, err
}

func (bc *BinanceClient) PlaceTradeTPOrder(ctx context.Context, t *models.Trade, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, *PreparedOrder, error) {
	bc.l.Printf("executing take-profit order for trade ID: %d", t.ID)
	po, err := bc.prepareTPOrder(ctx, t, f)
	if err != nil {
		bc.l.Printf("error during take-profit order preparation: %v", err)
		return nil, nil, err
	}

	res, err := bc.PlaceTPOrder(ctx, po)
	return res, po, err
}
