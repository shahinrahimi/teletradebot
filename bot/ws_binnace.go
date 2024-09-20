package bot

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

func (b *Bot) StartWebsocketServiceBinance(ctx context.Context) {
	go b.startUserDataStream724(ctx)
}

func (b *Bot) startUserDataStream724(ctx context.Context) {
	futures.UseTestnet = b.bc.UseTestnet
	for {
		// get new listenKey
		listenKey, err := b.bc.Client.NewStartUserStreamService().Do(ctx)
		if err != nil {
			b.l.Printf("Error starting user data stream: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		b.l.Printf("ListenKey acquired: %s", listenKey)

		// connect websocket
		doneC, stopC, err := futures.WsUserDataServe(listenKey, b.wsHandler, b.errHandler)
		if err != nil {
			b.l.Printf("Error connecting WebSocket: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		b.l.Println("WebSocket connection established. Awaiting events...")

		// Keep the connection alive by sending a ping every 30 minutes
		ticker := time.NewTicker(30 * time.Minute)

		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err := b.bc.Client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(ctx)
					if err != nil {
						b.l.Printf("Error keeping user data stream alive: %v", err)
					} else {
						b.l.Println("User data stream kept alive")
					}
				case <-ctx.Done():
					// Gracefully close the WebSocket connection when context is canceled
					close(stopC)
					b.l.Println("Context canceled, closing WebSocket connection...")
					return
				case <-doneC:
					return
				}
			}
		}()

		select {
		case <-doneC:
			close(stopC)
			b.l.Println("WebSocket connection closed, reconnecting...")
			time.Sleep(5 * time.Second)
		case <-ctx.Done():
			// Handle context cancellation in the main loop
			close(stopC)
			b.l.Println("Context canceled, exiting WebSocket loop...")
			return
		}
	}
}
