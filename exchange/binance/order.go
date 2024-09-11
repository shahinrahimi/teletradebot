package binance

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

var (
	TradeDescribers = map[int]*TradeDescriber{}
)

type PreparedOrder struct {
	Symbol     string
	Quantity   string
	StopPrice  string
	Side       futures.SideType
	Expiration time.Duration
}

type TradeDescriber struct {
	From  string
	Till  string
	Open  string
	High  string
	Low   string
	Close string
	SP    string // strop price or entry
	TP    string // take-profit price
	SL    string // take-loss price
}

func (td *TradeDescriber) ToTelegramString(t *models.Trade) string {
	sizeString := fmt.Sprintf("%.1f%%", float64(t.Size))
	slString := fmt.Sprintf("%.1f%%", float64((t.StopLoss - 100)))
	tpString := fmt.Sprintf("%.1f%%", float64((t.TakeProfit - 100)))

	msg := fmt.Sprintf("Trade ID %d\n\n", t.ID)
	msg = fmt.Sprintf("%s From:  %s\n Till:  %s\n Open:  %s\n High:  %s\n Low:  %s\n Close:  %s\n\n", msg, td.From, td.Till, td.Open, td.High, td.Low, td.Close)
	msg = fmt.Sprintf("%sTrading:\n", msg)
	msg = fmt.Sprintf("%s Entry %s at %s with %s of balance.\n", msg, t.Side, td.SP, sizeString)
	msg = fmt.Sprintf("%s TP at %s with %s.\n", msg, td.TP, tpString)
	msg = fmt.Sprintf("%s SL at %s with %s.\n", msg, td.SL, slString)
	return msg
}

func (bc *BinanceClient) getTradeLatestDescriber(ctx context.Context, t *models.Trade) (*TradeDescriber, error) {
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	if err != nil {
		return nil, err
	}
	sp, err := bc.calculateStopPrice(t, k, s)
	if err != nil {
		return nil, err
	}
	sl, err := bc.calculateStopLossPrice(t, k, s, sp)
	if err != nil {
		return nil, err
	}
	tp, err := bc.calculateTakeProfitPrice(t, k, s, sp)
	if err != nil {
		return nil, err
	}

	from := utils.ConvertTime(k.OpenTime).Format("2006-01-02 15:04:05")
	till := utils.ConvertTime(k.CloseTime).Add(time.Second).Format("2006-01-02 15:04:05")

	return &TradeDescriber{
		From:  from,
		Till:  till,
		Open:  k.Open,
		Close: k.Close,
		High:  k.High,
		Low:   k.Low,
		SP:    sp,
		TP:    tp,
		SL:    sl,
	}, nil
}

func (bc *BinanceClient) GetTradeDescriber(ctx context.Context, t *models.Trade) (*TradeDescriber, error) {
	td, exist := TradeDescribers[t.ID]
	if exist {
		return td, nil
	}
	return bc.getTradeLatestDescriber(ctx, t)
}

func (bc *BinanceClient) PrepareOrder(ctx context.Context, t *models.Trade) (*PreparedOrder, error) {
	var po PreparedOrder
	q, err := bc.getQuantity(t)
	if err != nil {
		return nil, err
	}
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	if err != nil {
		return nil, err
	}
	stopPrice, err := bc.calculateStopPrice(t, k, s)
	if err != nil {
		return nil, err
	}
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeBuy
	} else {
		side = futures.SideTypeSell
	}
	candleDuration, err := types.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}

	candleCloseTime := utils.ConvertTime(k.CloseTime)
	remainingTime := candleDuration + time.Until(candleCloseTime)
	if remainingTime < 0 {
		return nil, fmt.Errorf("remaining time should not be negative number: %d", remainingTime)
	}

	po.Symbol = t.Symbol
	po.Quantity = q
	po.Side = side
	po.StopPrice = stopPrice
	po.Expiration = remainingTime

	return &po, nil
}

func (bc *BinanceClient) PlacePreparedOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		StopPrice(p.StopPrice).
		Type(futures.OrderTypeStopMarket).
		WorkingType(futures.WorkingTypeMarkPrice)

	return order.Do(ctx)
}

func (bc *BinanceClient) CancelOrder(ctx context.Context, orderID int64, symbol string) (*futures.CancelOrderResponse, error) {
	order := bc.Client.NewCancelOrderService().
		OrderID(orderID).
		Symbol(symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) GetOrder(ctx context.Context, orderID int64, symbol string) (*futures.Order, error) {
	order := bc.Client.NewGetOrderService().
		OrderID(orderID).
		Symbol(symbol)
	return order.Do(ctx)
}

func (bc *BinanceClient) PrepareStopLossOrder(ctx context.Context, t *models.Trade, ou *futures.WsOrderTradeUpdate) (*PreparedOrder, error) {
	var po PreparedOrder
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	if err != nil {
		return nil, err
	}
	// TODO check if activation Price is true, or need to change it to other price like average price
	stopPrice, err := bc.calculateStopLossPrice(t, k, s, ou.StopPrice)
	if err != nil {
		return nil, err
	}
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}

	po.Symbol = t.Symbol
	po.Quantity = ou.OriginalQty
	po.Side = side
	po.StopPrice = stopPrice
	po.Expiration = 0

	return &po, nil
}

func (bc *BinanceClient) PrepareTakeProfitOrder(ctx context.Context, t *models.Trade, ou *futures.WsOrderTradeUpdate) (*PreparedOrder, error) {
	var po PreparedOrder
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	if err != nil {
		return nil, err
	}
	// TODO check if activation Price is true, or need to change it to other price like average price
	stopPrice, err := bc.calculateTakeProfitPrice(t, k, s, ou.StopPrice)
	if err != nil {
		return nil, err
	}
	var side futures.SideType
	if t.Side == types.SIDE_L {
		side = futures.SideTypeSell
	} else {
		side = futures.SideTypeBuy
	}

	po.Symbol = t.Symbol
	po.Quantity = ou.OriginalQty
	po.Side = side
	po.StopPrice = stopPrice
	po.Expiration = 0

	return &po, nil
}

func (bc *BinanceClient) PlacePreparedStopLossOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		WorkingType(futures.WorkingTypeMarkPrice).
		Type(futures.OrderTypeStopMarket).
		StopPrice(p.StopPrice)
	return order.Do(ctx)
}

func (bc *BinanceClient) PlacePreparedTakeProfitOrder(ctx context.Context, p *PreparedOrder) (*futures.CreateOrderResponse, error) {
	order := bc.Client.NewCreateOrderService().
		Symbol(p.Symbol).
		Side(p.Side).
		Quantity(p.Quantity).
		WorkingType(futures.WorkingTypeMarkPrice).
		Type(futures.OrderTypeTakeProfitMarket).
		StopPrice(p.StopPrice)
	return order.Do(ctx)
}
