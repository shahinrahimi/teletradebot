package binance

import (
	"context"
	"time"
)

func (bc *BinanceClient) StartPolling(ctx context.Context) {
	// Weight: 1
	// info like symbol price and symbol availability
	go bc.pollExchangeInfo(ctx, time.Hour)
}

func (bc *BinanceClient) pollExchangeInfo(ctx context.Context, interval time.Duration) {
	for {
		res, err := bc.Client.NewExchangeInfoService().Do(ctx)
		if err != nil {
			bc.l.Printf("Error fetching exchange info: %v", err)
			continue
		}
		bc.lastExchangeInfo = res
		time.Sleep(interval)
	}
}
