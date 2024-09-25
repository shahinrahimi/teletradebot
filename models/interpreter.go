package models

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

type Interpreter struct {
	Balance         float64
	Price           float64
	Quantity        float64
	ReverseQuantity float64
	// these fields copy from trade
	TradeID           int64
	Symbol            string
	Side              string
	Size              int
	StopLossSize      int
	TakeProfitSize    int
	ReverseMultiplier int
	TimeFrameDur      time.Duration

	// these fields fetched from exchange
	OpenTime  time.Time
	CloseTime time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	// these fields calculated
	EntryPrice             float64
	StopLossPrice          float64
	TakeProfitPrice        float64
	ReverseEntryPrice      float64
	ReverseStopLossPrice   float64
	ReverseTakeProfitPrice float64

	PricePrecision    int // use in binance exchange
	QuantityPrecision int // use in binance exchange

	TickSize    float64 // use in bitmex exchange
	LotSize     float64 // use in bitmex exchange
	MaxOrderQty float64 // use in bitmex exchange
}

func (i *Interpreter) AdjustPriceForBinance(value float64) string {
	return fmt.Sprintf("%.*f", i.PricePrecision, value)
}

func (i *Interpreter) AdjustQuantityForBinance(value float64) string {
	return fmt.Sprintf("%.*f", i.QuantityPrecision, value)
}

func (i *Interpreter) AdjustPriceForBitmex(value float64) float64 {
	return math.Round(value/i.TickSize) * i.TickSize
}

func (i *Interpreter) AdjustQuantityForBitmex(value float64) float64 {
	return math.Floor(value/i.LotSize) * i.LotSize
}

func (i *Interpreter) getPriceString(price float64) string {
	var priceStr string
	if i.PricePrecision > 0 {
		priceStr = fmt.Sprintf("%.*f", i.PricePrecision, price)
	} else if i.TickSize > 0 {
		p := math.Floor(price/i.TickSize) * i.TickSize
		precision := int(math.Abs(math.Log10(i.TickSize)))
		priceStr = fmt.Sprintf("%.*f", precision, p)
	} else {
		priceStr = fmt.Sprintf("%f", price)
	}
	return priceStr
}

func (i *Interpreter) CalculateExpiration() time.Duration {
	return i.TimeFrameDur + time.Until(i.CloseTime)
}
func (i *Interpreter) Describe() string {

	format := "2006-01-02 15:04:05"
	from := i.OpenTime.Local().Format(format)
	till := i.CloseTime.Local().Format(format)

	size := fmt.Sprintf("%.1f%%", float64(i.Size))
	reverseSize := fmt.Sprintf("%.1f%%", float64(i.ReverseMultiplier)*float64(i.Size))
	stopLossSize := fmt.Sprintf("%.1f%%", float64((i.StopLossSize - 100)))
	takeProfitSize := fmt.Sprintf("%.1f%%", float64((i.TakeProfitSize - 100)))

	open := i.getPriceString(i.Open)
	close := i.getPriceString(i.Close)
	high := i.getPriceString(i.High)
	low := i.getPriceString(i.Low)

	entryPrice := i.getPriceString(i.EntryPrice)
	stopLossPrice := i.getPriceString(i.StopLossPrice)
	takeProfitPrice := i.getPriceString(i.TakeProfitPrice)
	reverseStopLossPrice := i.getPriceString(i.ReverseStopLossPrice)
	reverseTakeProfitPrice := i.getPriceString(i.ReverseTakeProfitPrice)

	var reverseSide string
	if i.Side == types.SIDE_L {
		reverseSide = types.SIDE_S
	} else {
		reverseSide = types.SIDE_L
	}

	var expiration string
	if i.CalculateExpiration() > 0 {
		expiration = utils.FriendlyDuration(i.CalculateExpiration())
	} else {
		expiration = "∞"
	}

	msg := ""

	msg = fmt.Sprintf("Trade ID %d\n\n", i.TradeID)
	msg = fmt.Sprintf("%sFrom:  %s\nTill:  %s\nOpen:  %s\nHigh:  %s\nLow:  %s\nClose:  %s\n\n", msg, from, till, open, high, low, close)
	msg = fmt.Sprintf("%sTrading:\n", msg)
	msg = fmt.Sprintf("%sEntry %s at %s with %s of balance.\n", msg, i.Side, entryPrice, size)
	msg = fmt.Sprintf("%sTP at %s with %s.\n", msg, takeProfitPrice, takeProfitSize)
	msg = fmt.Sprintf("%sSL at %s with %s.\n\n", msg, stopLossPrice, stopLossSize)
	if i.ReverseMultiplier > 0 {
		msg = fmt.Sprintf("%sReverse:\n", msg)
		msg = fmt.Sprintf("%sEntry %s at %s with %s of balance.\n", msg, reverseSide, stopLossPrice, reverseSize)
		msg = fmt.Sprintf("%sTP at %s with %s.\n", msg, reverseTakeProfitPrice, takeProfitSize)
		msg = fmt.Sprintf("%sSL at %s with %s.\n\n", msg, reverseStopLossPrice, stopLossSize)
	}
	msg = fmt.Sprintf("%sExpiration: %s", msg, expiration)
	return msg
}

func (i *Interpreter) getSideBinance() futures.SideType {
	if i.Side == types.SIDE_S {
		return futures.SideTypeSell
	}
	return futures.SideTypeBuy
}

func (i *Interpreter) getOppositeSideBinance() futures.SideType {
	if i.Side == types.SIDE_S {
		return futures.SideTypeBuy
	}
	return futures.SideTypeSell
}

func (i *Interpreter) getSideBitmex() swagger.SideType {
	if i.Side == types.SIDE_S {
		return swagger.SideTypeSell
	}
	return swagger.SideTypeBuy
}

func (i *Interpreter) getOppositeSideBitmex() swagger.SideType {
	if i.Side == types.SIDE_S {
		return swagger.SideTypeBuy
	}
	return swagger.SideTypeSell
}

func (i *Interpreter) GetOrderExecutionBinance(executionType types.ExecutionType, orderIDStr string) *OrderExecutionBinance {
	var oeb *OrderExecutionBinance
	switch executionType {
	case types.GetOrderExecution, types.CancelOrderExecution:
		orderID, err := utils.ConvertOrderIDtoBinanceOrderID(orderIDStr)
		if err != nil {
			log.Panicf("invalid order id: %s", orderIDStr)
		}
		oeb = &OrderExecutionBinance{
			OrderID: orderID,
			Symbol:  i.Symbol,
		}
	case types.StopPriceExecution:
		side := i.getSideBinance()
		p := fmt.Sprintf("%.*f", i.PricePrecision, i.EntryPrice)
		q := fmt.Sprintf("%.*f", i.QuantityPrecision, i.Quantity)
		oeb = &OrderExecutionBinance{
			Symbol:    i.Symbol,
			Side:      side,
			StopPrice: p,
			Quantity:  q,
		}
	case types.TakeProfitExecution:
		side := i.getOppositeSideBinance()
		p := fmt.Sprintf("%.*f", i.PricePrecision, i.TakeProfitPrice)
		q := fmt.Sprintf("%.*f", i.QuantityPrecision, i.Quantity)
		oeb = &OrderExecutionBinance{
			Symbol:    i.Symbol,
			Side:      side,
			StopPrice: p,
			Quantity:  q,
		}
	case types.StopLossExecution:
		side := i.getOppositeSideBinance()
		p := fmt.Sprintf("%.*f", i.PricePrecision, i.StopLossPrice)
		q := fmt.Sprintf("%.*f", i.QuantityPrecision, i.Quantity)
		oeb = &OrderExecutionBinance{
			Symbol:    i.Symbol,
			Side:      side,
			StopPrice: p,
			Quantity:  q,
		}
	case types.ReverseStopPriceExecution:
		side := i.getOppositeSideBinance()
		p := fmt.Sprintf("%.*f", i.PricePrecision, i.ReverseEntryPrice)
		q := fmt.Sprintf("%.*f", i.QuantityPrecision, i.ReverseQuantity)
		oeb = &OrderExecutionBinance{
			Symbol:    i.Symbol,
			Side:      side,
			StopPrice: p,
			Quantity:  q,
		}
	case types.ReverseStopLossExecution:
		side := i.getSideBinance()
		p := fmt.Sprintf("%.*f", i.PricePrecision, i.ReverseStopLossPrice)
		q := fmt.Sprintf("%.*f", i.QuantityPrecision, i.ReverseQuantity)
		oeb = &OrderExecutionBinance{
			Symbol:    i.Symbol,
			Side:      side,
			StopPrice: p,
			Quantity:  q,
		}
	case types.ReverseTakeProfitExecution:
		side := i.getSideBinance()
		p := fmt.Sprintf("%.*f", i.PricePrecision, i.ReverseTakeProfitPrice)
		q := fmt.Sprintf("%.*f", i.QuantityPrecision, i.ReverseQuantity)
		oeb = &OrderExecutionBinance{
			Symbol:    i.Symbol,
			Side:      side,
			StopPrice: p,
			Quantity:  q,
		}
	case types.ClosePositionExecution:
		side := i.getOppositeSideBinance()
		q := fmt.Sprintf("%.*f", i.QuantityPrecision, i.Quantity)
		oeb = &OrderExecutionBinance{
			Symbol:   i.Symbol,
			Side:     side,
			Quantity: q,
		}
	case types.ClosePositionReverseExecution:
		side := i.getSideBinance()
		q := fmt.Sprintf("%.*f", i.QuantityPrecision, i.ReverseQuantity)
		oeb = &OrderExecutionBinance{
			Symbol:   i.Symbol,
			Side:     side,
			Quantity: q,
		}
	default:
		log.Panicf("invalid execution type: %s", executionType)
	}
	return oeb
}
func (i *Interpreter) GetOrderExecutionBitmex(ExecutionType types.ExecutionType, orderIDStr string) *OrderExecutionBitmex {
	var oeb *OrderExecutionBitmex
	switch ExecutionType {
	case types.GetOrderExecution, types.CancelOrderExecution:
		oeb = &OrderExecutionBitmex{
			OrderID: orderIDStr,
			Symbol:  i.Symbol,
		}
	case types.StopPriceExecution:
		side := i.getSideBitmex()
		p := math.Round(i.EntryPrice/i.TickSize) * i.TickSize
		q := math.Round(i.Quantity/i.LotSize) * i.LotSize
		oeb = &OrderExecutionBitmex{
			Symbol:    i.Symbol,
			Side:      side,
			Quantity:  q,
			StopPrice: p,
		}
	case types.TakeProfitExecution:
		side := i.getOppositeSideBitmex()
		p := math.Round(i.TakeProfitPrice/i.TickSize) * i.TickSize
		q := math.Round(i.Quantity/i.LotSize) * i.LotSize
		oeb = &OrderExecutionBitmex{
			Symbol:    i.Symbol,
			Side:      side,
			Quantity:  q,
			StopPrice: p,
		}
	case types.StopLossExecution:
		side := i.getOppositeSideBitmex()
		p := math.Round(i.StopLossPrice/i.TickSize) * i.TickSize
		q := math.Round(i.Quantity/i.LotSize) * i.LotSize
		oeb = &OrderExecutionBitmex{
			Symbol:    i.Symbol,
			Side:      side,
			Quantity:  q,
			StopPrice: p,
		}
	case types.ReverseStopPriceExecution:
		side := i.getOppositeSideBitmex()
		p := math.Round(i.ReverseEntryPrice/i.TickSize) * i.TickSize
		q := math.Round(i.ReverseQuantity/i.LotSize) * i.LotSize
		oeb = &OrderExecutionBitmex{
			Symbol:    i.Symbol,
			Side:      side,
			Quantity:  q,
			StopPrice: p,
		}
	case types.ReverseStopLossExecution:
		side := i.getSideBitmex()
		p := math.Round(i.ReverseStopLossPrice/i.TickSize) * i.TickSize
		q := math.Round(i.ReverseQuantity/i.LotSize) * i.LotSize
		oeb = &OrderExecutionBitmex{
			Symbol:    i.Symbol,
			Side:      side,
			Quantity:  q,
			StopPrice: p,
		}
	case types.ReverseTakeProfitExecution:
		side := i.getSideBitmex()
		p := math.Round(i.ReverseTakeProfitPrice/i.TickSize) * i.TickSize
		q := math.Round(i.ReverseQuantity/i.LotSize) * i.LotSize
		oeb = &OrderExecutionBitmex{
			Symbol:    i.Symbol,
			Side:      side,
			Quantity:  q,
			StopPrice: p,
		}
	case types.ClosePositionExecution, types.ClosePositionReverseExecution:
		oeb = &OrderExecutionBitmex{
			Symbol: i.Symbol,
		}
	default:
		log.Panicf("invalid execution type: %s", ExecutionType)
	}
	return oeb

}
