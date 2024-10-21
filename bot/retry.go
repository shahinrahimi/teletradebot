package bot

import (
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) retry(attempts int, delay time.Duration, t *models.Trade, f func() (interface{}, error)) (interface{}, error) {
	var err error
	for i := 0; i < attempts; i++ {
		b.DbgChan <- fmt.Sprintf("Attempting to perform action on trade, attempt number: %d, TradeID: %d", i+1, t.ID)
		res, err := f()
		if err != nil {
			b.DbgChan <- fmt.Sprintf("Failed to perform action on order, attempt number: %d, error: %v, TradeID: %d", i+1, err, t.ID)
			if apiErr, ok := err.(*common.APIError); ok {
				b.DbgChan <- fmt.Sprintf("Error code: %d, error message: %s", apiErr.Code, apiErr.Message)
				switch {
				case (apiErr.Code == -1007 || apiErr.Code == -1008):
					msg := fmt.Sprintf("Failed to perform action on order\nError: %s\nRetry after %s\n\nTrade ID: %d", apiErr.Message, utils.FriendlyDuration(delay), t.ID)
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
					msg := fmt.Sprintf("Failed to perform action on order\nError: %s\nRetry after %s\n\nTrade ID: %d", apiErr.Error(), utils.FriendlyDuration(delay), t.ID)
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
				b.DbgChan <- fmt.Sprintf("unexpected error happened in retrying function: %v", err)
				return nil, err
			}
		}
		return res, err
	}
	b.DbgChan <- fmt.Sprintf("All attempts failed to perform action on order, TradeID: %d", t.ID)
	return nil, err
}

func (b *Bot) retryDenyNotFound(attempts int, delay time.Duration, t *models.Trade, f func() (interface{}, error)) (interface{}, error) {
	var err error
	for i := 0; i < attempts; i++ {
		res, err := f()
		if err != nil {
			if apiErr, ok := err.(*common.APIError); ok {
				switch {
				case apiErr.Code == -2011:
					b.handleError(err, t.UserID, t.ID)
					// deny the error from binance api
					return nil, nil
				case (apiErr.Code == -1007 || apiErr.Code == -1008):
					msg := fmt.Sprintf("Failed to perform action on order\nError: %s\nRetry after %s\n\nTrade ID: %d", apiErr.Message, utils.FriendlyDuration(delay), t.ID)
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
				switch {
				case apiErr.Error() == "503 Service Unavailable":
					msg := fmt.Sprintf("Failed to perform action on order\nError: %s\nRetry after %s\n\nTrade ID: %d", apiErr.Error(), utils.FriendlyDuration(delay), t.ID)
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
				b.l.Printf("unexpected error happened in retrying function: %v", err)
				return nil, err
			}
		}
		return res, err
	}
	return nil, err
}
