package models

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/swagger"
)

type OrderExecutionBinance struct {
	OrderID   int64 // use for get and cancel the order
	Symbol    string
	Quantity  string
	StopPrice string
	Side      futures.SideType
}

type OrderExecutionBitmex struct {
	OrderID   string // use for get and cancel the order
	Symbol    string
	Quantity  float64
	StopPrice float64
	Side      swagger.SideType
}
