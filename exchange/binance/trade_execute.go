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
	bc.l.Printf("executing order for trade ID: %d", t.ID)
	po, err := bc.prepareOrder(ctx, t)
	if err != nil {
		bc.l.Printf("trade could not be executed, error in preparing state: %v", err)
		return
	}
	bc.l.Printf("Placing %s order with quantity %s and stop price %s expires in: %s", po.Side, po.Quantity, po.StopPrice, utils.FriendlyDuration(po.Expiration))
	tries := 0
	for {
		bc.l.Printf("tries: %d", tries)
		res, err := bc.placeOrder(ctx, po)
		if err != nil {
			bc.l.Printf("error while placing order: %v", err)
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
					bc.l.Printf("Im here with error 3: %v", err)
					bc.handleError(apiErr, t.UserID)
					return
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
			bc.l.Printf("error updating trade to DB: %v", err)
		}
		return
	}
}

func (bc *BinanceClient) executeSLOrder(t *models.Trade, f *futures.WsOrderTradeUpdate) {
	bc.l.Printf("executing stop-loss order for trade ID: %d", t.ID)
	po, err := bc.prepareSLOrder(context.Background(), t, f)
	if err != nil {
		bc.l.Printf("error during stop-loss order preparation: %v", err)
		return
	}
	var tries int
	for {
		bc.l.Printf("try %d to place stop-loss order for trade ID: %d", tries, t.ID)
		res, err := bc.placeSLOrder(context.Background(), po)
		if err != nil {
			if apiErr, ok := err.(*common.APIError); ok {
				bc.l.Printf("error placing stop-loss order: %v", err)
				switch {
				case (apiErr.Code == -1007 || apiErr.Code == -1008) && tries <= config.MaxTries:
					tries++
					bc.l.Printf("retrying after %d seconds", config.WaitForNextTries)
					msg := fmt.Sprintf("Failed to place a stop-loss order\nRetry after %d seconds ...\n\nTrade ID: %d", config.WaitForNextTries, t.ID)
					bc.MsgChan <- types.BotMessage{
						ChatID: t.UserID,
						MsgStr: msg,
					}
					time.Sleep(config.WaitForNextTries)
					continue
				default:
					bc.l.Printf("error placing stop-loss order: %v", err)
					bc.handleError(apiErr, t.UserID)
					return
				}
			} else {
				bc.l.Printf("unexpected error happened in placing a order: %v", err)
				return
			}
		}
		orderID := utils.ConvertBinanceOrderID(res.OrderID)
		bc.l.Printf("stop-loss order placed successfully for trade ID: %d, order ID: %s", t.ID, orderID)
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
	bc.l.Printf("executing take-profit order for trade ID: %d", t.ID)
	po, err := bc.prepareTPOrder(context.Background(), t, f)
	if err != nil {
		bc.l.Printf("error during take-profit order preparation: %v", err)
		return
	}
	bc.l.Printf("Placing %s take-profit order with quantity %s and stop price %s.", po.Side, po.Quantity, po.StopPrice)
	tries := 0
	for {
		bc.l.Printf("try %d to place take-profit order", tries)
		res, err := bc.placeTPOrder(context.Background(), po)
		if err != nil {
			if apiErr, ok := err.(*common.APIError); ok {
				bc.l.Printf("binance error: %v", apiErr)
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
					return
				}
			} else {
				bc.l.Printf("unexpected error happened in placing a order: %v", err)
				return
			}
		}
		orderID := utils.ConvertBinanceOrderID(res.OrderID)
		bc.l.Printf("take-profit order placed successfully for trade ID: %d, order ID: %s", t.ID, orderID)
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
