package binance

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/common"
)

var (
	initialBackoff = 5 * time.Second
	maxBackOff     = 5 * time.Minute
	backoff        = initialBackoff
)

func (bc *BinanceClient) StartPolling(ctx context.Context) {
	// Weight: 1
	// info like symbol price and symbol availability
	go bc.pollExchangeInfo(ctx, time.Hour)
}

func (bc *BinanceClient) pollExchangeInfo(ctx context.Context, interval time.Duration) {
	for {
		res, err := bc.client.NewExchangeInfoService().Do(ctx)
		if err != nil {
			bc.l.Printf("Error fetching exchange info: %v", err)
			if apiErr, ok := err.(*common.APIError); ok {
				if apiErr.Code == 429 {
					bc.l.Printf("hit binance rate limit: %+v", apiErr.Response)
					time.Sleep(backoff)
					if backoff > maxBackOff {
						backoff = maxBackOff
					} else {
						backoff *= 2
					}
					continue
				}
			} else {
				bc.l.Panicf("error fetching exchange info: %v", err)
			}
		}
		bc.lastExchangeInfo = res
		time.Sleep(interval)
	}
}
