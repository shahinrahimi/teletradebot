package binance

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) CancelTrade(ctx context.Context, t *models.Trade) (*futures.CancelOrderResponse, error) {
	orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.OrderID)
	if err != nil {
		bc.l.Printf("Unexpected issue: the trade's OrderID is not in a valid format for conversion: %v", err)
		return nil, err
	}
	return bc.cancelOrder(ctx, orderID, t.Symbol)
}
