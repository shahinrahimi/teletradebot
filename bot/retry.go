package bot

import (
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) retry(funcName string, denyNotFoundBinance bool, t *models.Trade, f func() (interface{}, error)) (interface{}, error) {
	attempts := config.MaxTries
	delay := config.WaitForNextTries
	var err error
	for i := 0; i < attempts; i++ {
		b.DbgChan <- fmt.Sprintf("Attempting to perform action(%s) on trade, attempt number: %d, TradeID: %d", funcName, i+1, t.ID)
		res, err := f()
		if err != nil {
			b.DbgChan <- fmt.Sprintf("Failed to perform action(%s) on trade, attempt number: %d, error: %v, TradeID: %d", funcName, i+1, err, t.ID)
			if apiErr, ok := err.(*common.APIError); ok {
				b.DbgChan <- fmt.Sprintf("Error code: %d, error message: %s", apiErr.Code, apiErr.Message)
				switch {
				case apiErr.Code == -2011 && denyNotFoundBinance:
					b.DbgChan <- fmt.Sprintf("Deny the error(%d) from binance api, TradeID: %d", apiErr.Code, t.ID)
					b.handleError(err, t.UserID, t.ID)
					// deny the error from binance api
					return nil, nil
				case (apiErr.Code == -1007 || apiErr.Code == -1008):
					msg := fmt.Sprintf("Failed to perform action(%s) on trade\nError: %s\nRetry after %s\n\nTrade ID: %d", funcName, apiErr.Message, utils.FriendlyDuration(delay), t.ID)
					b.MsgChan <- types.BotMessage{
						ChatID: t.UserID,
						MsgStr: msg,
					}
					time.Sleep(delay)
					continue
				default:
					return nil, err
				}
			} else if apiErr, ok := err.(swagger.GenericSwaggerError); ok {
				b.DbgChan <- fmt.Sprintf("Error message: %s", apiErr.Error())
				switch {
				case apiErr.Error() == "503 Service Unavailable":
					msg := fmt.Sprintf("Failed to perform action(%s) on trade\nError: %s\nRetry after %s\n\nTrade ID: %d", funcName, apiErr.Error(), utils.FriendlyDuration(delay), t.ID)
					b.MsgChan <- types.BotMessage{
						ChatID: t.UserID,
						MsgStr: msg,
					}
					time.Sleep(delay)
					continue
				default:
					return nil, err
				}
			} else if _, ok := err.(*types.BotError); ok {
				return nil, err
			} else {
				b.DbgChan <- fmt.Sprintf("unexpected error happened in retrying function for action(%s) error: %v", funcName, err)
				return nil, err
			}
		}
		return res, err
	}
	b.DbgChan <- fmt.Sprintf("All attempts failed to perform action on order, TradeID: %d", t.ID)
	return nil, err
}
