package binance

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
)

func (bc *BinanceClient) handleNewFilled(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	// freeze describer
	d, err := bc.FetchDescriber(context.Background(), t)
	if err != nil {
		bc.l.Printf("error fetching the describer %v", err)
	} else {
		models.SetDescriber(d, t.ID)
	}

	go bc.executeSLOrder(t, f)
	go bc.executeTPOrder(t, f)
}
