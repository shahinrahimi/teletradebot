package binance

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (bc *BinanceClient) handleNewFilled(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	td, err := bc.GetTradeDescriber(context.Background(), t)
	if err != nil {
		bc.l.Printf("error getting trade describer")
	}
	go bc.executeSLOrder(t, f)
	go bc.executeTPOrder(t, f)

	if td != nil {
		types.TradeDescribers[t.ID] = &types.TradeDescriber{
			From:  td.From,
			Till:  td.Till,
			Open:  td.Open,
			Close: td.Close,
			High:  td.High,
			Low:   td.Low,
			TP:    "otp.StopPrice",
			SL:    "otp.StopPrice",
			SP:    f.StopPrice,
		}
	}

}
