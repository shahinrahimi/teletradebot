package binance

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) executeTPOrder(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	po, err := bc.prepareTPOrder(context.Background(), t, f)
	if err != nil {
		bc.l.Printf("error during take-profit order preparation: %v", err)
		return
	}
	bc.l.Printf("Placing %s take-profit order with quantity %s and stop price %s.", po.Side, po.Quantity, po.StopPrice)
	tries := 0
	for {
		res, err := bc.placeTPOrder(context.Background(), po)
		if err != nil {
			if apiErr, ok := err.(*common.APIError); ok {
				switch {
				case (apiErr.Code == -1007 || apiErr.Code == -1008) && tries <= config.MaxTries:
					tries = tries + 1
					msg := fmt.Sprintf("Failed to place a take-profit order\nRetry after %d seconds ...\n\nTrade ID: %d", config.WaitForNextTries, t.ID)
					bc.MsgChan <- types.BotMessage{
						ChatID: t.UserID,
						MsgStr: msg,
					}
					time.Sleep(config.WaitForNextTries)
					continue
				default:
					bc.handleError(apiErr, t.UserID)
				}
			} else {
				bc.l.Printf("unexpected error happened in placing a order: %v", err)
				return
			}
		}
		orderID := utils.ConvertBinanceOrderID(res.OrderID)
		msg := fmt.Sprintf("take-profit order placed successfully.\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)

		bc.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}

		if err := bc.s.UpdateTradeTPOrder(t, orderID); err != nil {
			bc.l.Printf("error updating trade for TP orderID: %v", err)
		}
		return

	}
}
