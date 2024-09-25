package binance

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

func (bc *BinanceClient) StartWebsocketService(ctx context.Context, wsHandler func(event *futures.WsUserDataEvent), errHandler func(err error)) {
	go bc.startUserDataStream724(ctx, wsHandler, errHandler)
}

func (bc *BinanceClient) startUserDataStream724(ctx context.Context, wsHandler func(event *futures.WsUserDataEvent), errHandler func(err error)) {
	futures.UseTestnet = bc.UseTestnet
	for {
		// get new listenKey
		listenKey, err := bc.client.NewStartUserStreamService().Do(ctx)
		if err != nil {
			bc.l.Printf("Error starting user data stream: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		//b.l.Printf("ListenKey acquired: %s", listenKey)

		// connect websocket
		doneC, stopC, err := futures.WsUserDataServe(listenKey, wsHandler, errHandler)
		if err != nil {
			bc.l.Printf("Error connecting WebSocket: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		bc.l.Println("WebSocket connection established. Awaiting events...")

		// Keep the connection alive by sending a ping every 30 minutes
		ticker := time.NewTicker(30 * time.Minute)

		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err := bc.client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(ctx)
					if err != nil {
						bc.l.Printf("Error keeping user data stream alive: %v", err)
					} else {
						bc.l.Println("User data stream kept alive")
					}
				case <-ctx.Done():
					// Gracefully close the WebSocket connection when context is canceled
					close(stopC)
					bc.l.Println("Context canceled, closing WebSocket connection...")
					return
				case <-doneC:
					return
				}
			}
		}()

		select {
		case <-doneC:
			close(stopC)
			bc.l.Println("WebSocket connection closed, reconnecting...")
			time.Sleep(5 * time.Second)
		case <-ctx.Done():
			// Handle context cancellation in the main loop
			close(stopC)
			bc.l.Println("Context canceled, exiting WebSocket loop...")
			return
		}
	}
}
