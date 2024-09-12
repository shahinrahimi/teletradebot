package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
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

func (bc *BinanceClient) handleSLFilled(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.TPOrderID)
	if err != nil {
		bc.l.Panicf("Error converting take-profit order ID: %v", err)
		return
	}
	msg := fmt.Sprintf("🛑 Stop-loss order executed successfully.\n\nTrade ID: %d", t.ID)
	bc.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}
	if _, err := bc.cancelOrder(context.Background(), orderID, f.Symbol); err != nil {
		bc.l.Printf("Error canceling take-profit order: %v", err)
		return
	}
	msg = fmt.Sprintf("Take-profit order has been canceled.\n\nTrade ID: %d", t.ID)
	bc.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}
}

func (bc *BinanceClient) handleTPFilled(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.SLOrderID)
	if err != nil {
		bc.l.Panicf("Error converting stop-loss order ID: %v", err)
		return
	}
	msg := fmt.Sprintf("✅ Take-profit order executed successfully.\n\nTrade ID: %d", t.ID)
	bc.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}
	if _, err := bc.cancelOrder(context.Background(), orderID, f.Symbol); err != nil {
		bc.l.Printf("Error canceling stop-loss order: %v", err)
		return
	}
	msg = fmt.Sprintf("Stop-loss order has been canceled.\n\nTrade ID: %d", t.ID)
	bc.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}
}