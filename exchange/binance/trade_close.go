package binance

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) CloseTrade(ctx context.Context, t *models.Trade) {
	var orderID int64
	orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.OrderID)
	if err != nil {
		bc.l.Printf("error converting orderID to binance orderID: %v", err)
		return
	}

	o, err := bc.getOrder(ctx, orderID, t.Symbol)
	if err != nil {
		bc.l.Printf("error getting order: %v", err)
	}
	var side futures.SideType
	if o.Side == futures.SideTypeBuy {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}
	switch {
	case o.Status == futures.OrderStatusTypeCanceled:
		bc.l.Printf("the state is : %v", o.Status)
	case o.Status == futures.OrderStatusTypeExpired:
		bc.l.Printf("the state is : %v", o.Status)
	case o.Status == futures.OrderStatusTypeFilled:
		go bc.closeOrder(ctx, o.OrigQuantity, side, o.Symbol)
	case o.Status == futures.OrderStatusTypeNew:
		go bc.cancelOrder(ctx, o.OrderID, o.Symbol)
	case o.Status == futures.OrderStatusTypePartiallyFilled:
		go bc.closeOrder(ctx, o.ExecutedQuantity, side, o.Symbol)
	default:
		bc.l.Printf("the state is : %v", o.Status)
	}

	orderID, err = utils.ConvertOrderIDtoBinanceOrderID(t.TPOrderID)
	if err != nil {
		bc.l.Printf("error getting order: %v", err)
	} else {
		go bc.cancelOrder(ctx, orderID, t.Symbol)
	}

	orderID, err = utils.ConvertOrderIDtoBinanceOrderID(t.SLOrderID)
	if err != nil {
		bc.l.Printf("error getting order: %v", err)
	} else {
		go bc.cancelOrder(ctx, orderID, t.Symbol)
	}

}
