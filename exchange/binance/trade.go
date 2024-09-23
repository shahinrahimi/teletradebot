package binance

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
)

func (bc *BinanceClient) PlaceTrade(ctx context.Context, t *models.Trade) (*futures.CreateOrderResponse, *models.Describer, error) {
	bc.l.Printf("executing order for trade ID: %d", t.ID)
	// fetch describer
	d, err := bc.FetchDescriber(ctx, t)
	if err != nil {
		bc.l.Printf("error fetching the describer %v", err)
		return nil, nil, err
	}
	// prepare order
	po, err := bc.prepareMainOrder(ctx, d)
	if err != nil {
		bc.l.Printf("trade could not be executed, error in preparing state: %v", err)
		return nil, nil, err
	}
	res, err := bc.PlaceOrder(ctx, po)
	return res, d, err
}

func (bc *BinanceClient) PlaceTradeSLOrder(ctx context.Context, t *models.Trade, d *models.Describer, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	bc.l.Printf("executing stop-loss order for trade ID: %d", t.ID)
	po := bc.prepareStopLossOrder(ctx, d, f)
	res, err := bc.PlaceSLOrder(ctx, po)
	return res, err
}

func (bc *BinanceClient) PlaceTradeTPOrder(ctx context.Context, t *models.Trade, d *models.Describer, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	bc.l.Printf("executing take-profit order for trade ID: %d", t.ID)
	po := bc.prepareTakeProfitOrder(ctx, d, f)
	res, err := bc.PlaceTPOrder(ctx, po)
	return res, err
}

// func (bc *BinanceClient) PlaceReverseTradeSLOrder(ctx context.Context, t *models.Trade, d *models.Describer, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
// 	bc.l.Printf("executing reverse stop-loss order for trade ID: %d", t.ID)
// 	po := bc.prepareDescriberForReverseStopLossOrder(ctx, d, t, f)
// 	res, err := bc.PlaceSLOrder(ctx, po)
// 	return res, err
// }

// func (bc *BinanceClient) PlaceReverseTradeTPOrder(ctx context.Context, t *models.Trade, d *models.Describer, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
// 	bc.l.Printf("executing reverse take-profit order for trade ID: %d", t.ID)
// 	po := bc.prepareDescriberForReverseTakeProfitOrder(ctx, d, t, f)
// 	res, err := bc.PlaceTPOrder(ctx, po)
// 	return res, err
// }
