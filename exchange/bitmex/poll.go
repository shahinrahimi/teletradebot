package bitmex

import (
	"context"
	"time"
)

func (mc *BitmexClient) StartPolling(ctx context.Context) {
	//go mc.pollInstruments(ctx, time.Minute)
	//go mc.pollMargin(ctx, time.Minute)
}

func (mc *BitmexClient) pollInstruments(ctx context.Context, interval time.Duration) {
	for {
		instruments, _, err := mc.client.InstrumentApi.InstrumentGet(mc.auth, nil)
		if err != nil {
			mc.l.Panicf("error fetching instruments: %v", err)
			continue
		}
		// Access the rate limit headers
		// rateLimitRemaining := res.Header.Get("X-RateLimit-Remaining")
		// rateLimitReset := res.Header.Get("X-RateLimit-Reset")
		// mc.l.Printf("Remaining requests: %s, Resets in: %s seconds\n", rateLimitRemaining, rateLimitReset)
		mc.l.Printf("instruments length: %d", len(instruments))
		time.Sleep(interval)

	}
}

func (mc *BitmexClient) pollMargin(ctx context.Context, interval time.Duration) {
	for {
		_, _, err := mc.client.UserApi.UserGetMargin(mc.auth, nil)
		if err != nil {
			mc.l.Panicf("error fetching instruments: %v", err)
			continue
		}
		// Access the rate limit headers
		// rateLimitRemaining := res.Header.Get("X-RateLimit-Remaining")
		// rateLimitReset := res.Header.Get("X-RateLimit-Reset")
		// mc.l.Printf("Remaining requests: %s, Resets in: %s seconds\n", rateLimitRemaining, rateLimitReset)
		time.Sleep(interval)

	}
}
