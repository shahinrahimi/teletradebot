package binanceold

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

func (bc *BinanceClient) PlaceTradeSLOrder(ctx context.Context, d *models.Describer, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	bc.l.Printf("executing stop-loss order for trade ID: %d", d.TradeID)
	po := bc.prepareStopLossOrder(d, f)
	res, err := bc.PlaceSLOrder(ctx, po)
	return res, err
}

func (bc *BinanceClient) PlaceTradeTPOrder(ctx context.Context, d *models.Describer, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	bc.l.Printf("executing take-profit order for trade ID: %d", d.TradeID)
	po := bc.prepareTakeProfitOrder(d, f)
	res, err := bc.PlaceTPOrder(ctx, po)
	return res, err
}

func (bc *BinanceClient) PlaceTradeReverseMainOrder(ctx context.Context, d *models.Describer, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	bc.l.Printf("executing reverse main order for trade ID: %d", d.TradeID)
	po, err := bc.prepareReverseMainOrder(d, f)
	if err != nil {
		bc.l.Printf("error preparing reverse main order: %v", err)
		return nil, err
	}
	res, err := bc.PlaceOrder(ctx, po)
	return res, err
}

func (bc *BinanceClient) PlaceTradeReverseStopLossOrder(ctx context.Context, d *models.Describer, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	bc.l.Printf("executing reverse stop-loss order for trade ID: %d", d.TradeID)
	po := bc.prepareReverseStopLossOrder(d, f)
	res, err := bc.PlaceSLOrder(ctx, po)
	return res, err
}

func (bc *BinanceClient) PlaceTradeReverseTakeProfitOrder(ctx context.Context, d *models.Describer, f *futures.WsOrderTradeUpdate) (*futures.CreateOrderResponse, error) {
	bc.l.Printf("executing reverse take-profit order for trade ID: %d", d.TradeID)
	po := bc.prepareReverseTakeProfitOrder(d, f)
	res, err := bc.PlaceTPOrder(ctx, po)
	return res, err
}
