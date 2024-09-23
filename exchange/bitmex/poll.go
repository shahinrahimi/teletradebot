package bitmex

import (
	"context"
	"time"
)

func (mc *BitmexClient) StartPolling(ctx context.Context) {
	// Weight: 1
	// info like symbol price and symbol availability
	go mc.pollActiveInstruments(ctx, time.Hour)
}

func (mc *BitmexClient) pollActiveInstruments(ctx context.Context, interval time.Duration) {
	for {
		instruments, _, err := mc.client.InstrumentApi.InstrumentGetActive(mc.auth)
		if err != nil {
			mc.l.Printf("Error fetching exchange info: %v", err)
			continue
		}
		mc.activeInstruments = instruments
		time.Sleep(interval)
	}
}
