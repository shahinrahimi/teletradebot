package bitmex

import (
	"context"
	"time"

	"github.com/shahinrahimi/teletradebot/swagger"
)

var (
	initialBackoff = 5 * time.Second
	maxBackOff     = 5 * time.Minute
	backoff        = initialBackoff
)

func (mc *BitmexClient) StartPolling(ctx context.Context) {
	// Weight: 1
	// info like symbol price and symbol availability
	go mc.pollActiveInstruments(ctx, time.Hour)
}

func (mc *BitmexClient) pollActiveInstruments(ctx context.Context, interval time.Duration) {
	for {
		ctx = mc.getAuthContext(ctx)
		instruments, _, err := mc.client.InstrumentApi.InstrumentGetActive(ctx)
		if err != nil {
			mc.l.Printf("Error fetching exchange info: %v", err)
			if apiErr, ok := err.(swagger.GenericSwaggerError); ok {
				if apiErr.Error() == "429 Too Many Requests" {
					mc.l.Printf("hit bitmex rate limit: %+v", apiErr.Body())
					time.Sleep(backoff)
					if backoff > maxBackOff {
						backoff = maxBackOff
					} else {
						backoff *= 2
					}
					continue
				}
			} else {
				mc.l.Panicf("error fetching exchange info: %v", err)
			}
		}
		mc.activeInstruments = instruments
		time.Sleep(interval)
	}
}
