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

func (bc *BinanceClient) ExecuteTrade(ctx context.Context, t *models.Trade, isReplacement bool) {
	po, err := bc.prepareOrder(ctx, t)
	if err != nil {
		bc.l.Printf("trade could not be executed, error in preparing state: %v", err)
		return
	}
	bc.l.Printf("Placing %s order with quantity %s and stop price %s expires in: %s", po.Side, po.Quantity, po.StopPrice, utils.FriendlyDuration(po.Expiration))
	tries := 0
	for {
		res, err := bc.placeOrder(ctx, po)
		if err != nil {
			if apiErr, ok := err.(*common.APIError); ok {
				switch {
				case (apiErr.Code == -1007 || apiErr.Code == -1008) && tries <= config.MaxTries:
					tries = tries + 1
					msg := fmt.Sprintf("Failed to place a order\nRetry after %d seconds ...\n\nTrade ID: %d", config.WaitForNextTries, t.ID)
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

		// schedule for replacement
		go bc.scheduleOrderReplacement(ctx, po.Expiration, res.OrderID, t)
		orderID := utils.ConvertBinanceOrderID(res.OrderID)

		// message the user
		var msg string
		if isReplacement {
			msg = fmt.Sprintf("Order replaced successfully\n\nNewOrder ID: %s\nTrade ID: %d", orderID, t.ID)
		} else {
			msg = fmt.Sprintf("Order placed successfully\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)
		}
		bc.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
		// update trade state
		if err := bc.s.UpdateTradePlaced(t, orderID); err != nil {
			bc.l.Printf("error updating trade to DB")
		}
		return
	}
}

func (bc *BinanceClient) executeSLOrder(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	po, err := bc.prepareSLOrder(context.Background(), t, f)
	if err != nil {
		bc.l.Printf("error during stop-loss order preparation: %v", err)
		return
	}
	bc.l.Printf("Placing %s stop-loos order with quantity %s and stop price %s.", po.Side, po.Quantity, po.StopPrice)
	tries := 0
	for {
		res, err := bc.placeSLOrder(context.Background(), po)
		if err != nil {
			if apiErr, ok := err.(*common.APIError); ok {
				switch {
				case (apiErr.Code == -1007 || apiErr.Code == -1008) && tries <= config.MaxTries:
					tries = tries + 1
					msg := fmt.Sprintf("Failed to place a stop-loss order\nRetry after %d seconds ...\n\nTrade ID: %d", config.WaitForNextTries, t.ID)
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
		msg := fmt.Sprintf("Stop-loss order placed successfully.\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)

		bc.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}

		if err := bc.s.UpdateTradeSLOrder(t, orderID); err != nil {
			bc.l.Printf("error updating trade for SL orderID: %v", err)
		}
		return

	}
}

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
