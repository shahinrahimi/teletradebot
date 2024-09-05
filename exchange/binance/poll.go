package binance

import (
	"context"
	"time"
)

func (bc *BinanceClient) StartPolling(ctx context.Context) {
	// Weight: ?
	//go bc.pollListenKey(ctx, time.Hour)
	// Weight: 5
	go bc.pollAccount(ctx, time.Minute)
	// Weight: 1
	go bc.pollExchangeInfo(ctx, time.Minute)
	// Weight: 2
	go bc.pollSymbolPrices(ctx, time.Minute)
}

func (bc *BinanceClient) pollSymbolPrices(ctx context.Context, interval time.Duration) {
	for {
		res, err := bc.client.NewListPricesService().Do(ctx)
		if err != nil {
			bc.l.Printf("error fetching symbol prices: %v", err)
			continue
		}
		bc.lastSymbolPrices = res
		time.Sleep(interval)
	}
}

func (bc *BinanceClient) pollExchangeInfo(ctx context.Context, interval time.Duration) {
	for {
		res, err := bc.client.NewExchangeInfoService().Do(ctx)
		if err != nil {
			bc.l.Printf("Error fetching exchange info: %v", err)
			continue
		}
		bc.lastExchangeInfo = res
		time.Sleep(interval)
	}
}

func (bc *BinanceClient) pollAccount(ctx context.Context, interval time.Duration) {
	for {
		res, err := bc.client.NewGetAccountService().Do(ctx)
		if err != nil {
			bc.l.Printf("Error fetching exchange info: %v", err)
			continue
		}
		bc.lastAccount = res
		time.Sleep(interval)
	}
}

func (bc *BinanceClient) pollListenKey(ctx context.Context, interval time.Duration) {
	for {
		listenKey, err := bc.client.NewStartUserStreamService().Do(ctx)
		if err != nil {
			bc.l.Printf("Error fetching listen key: %v", err)
			continue
		}
		bc.ListenKey = listenKey
		time.Sleep(interval)

	}
}
