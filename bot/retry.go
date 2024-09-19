package bot

import (
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) retry(attempts int, delay time.Duration, t *models.Trade, f func() (interface{}, interface{}, error)) (interface{}, interface{}, error) {
	var err error
	for i := 0; i < attempts; i++ {
		res, po, err := f()
		if err != nil {
			if apiErr, ok := err.(*common.APIError); ok {
				switch {
				case (apiErr.Code == -1007 || apiErr.Code == -1008):
					// TODO add action string so the message have meaning (not place order)
					msg := fmt.Sprintf("Failed to place a order\nRetry after %s ...\n\nTrade ID: %d", utils.FriendlyDuration(delay), t.ID)
					b.MsgChan <- types.BotMessage{
						ChatID: t.UserID,
						MsgStr: msg,
					}
					time.Sleep(delay)
					continue
				default:
					return nil, nil, err
				}
			} else {
				b.l.Printf("unexpected error happened in placing a order: %v", err)
				return nil, nil, err
			}
		}
		return res, po, err
	}
	return nil, nil, err
}

func (b *Bot) retry2(attempts int, delay time.Duration, t *models.Trade, f func() (interface{}, error)) (interface{}, error) {
	var err error
	for i := 0; i < attempts; i++ {
		res, err := f()
		if err != nil {
			if apiErr, ok := err.(*common.APIError); ok {
				switch {
				case (apiErr.Code == -1007 || apiErr.Code == -1008):
					msg := fmt.Sprintf("Failed to place a order\nRetry after %s ...\n\nTrade ID: %d", utils.FriendlyDuration(delay), t.ID)
					b.MsgChan <- types.BotMessage{
						ChatID: t.UserID,
						MsgStr: msg,
					}
					time.Sleep(delay)
					continue
				default:
					return nil, err
				}
			} else {
				b.l.Printf("unexpected error happened in placing a order: %v", err)
				return nil, err
			}
		}
		return res, err
	}
	return nil, err
}
