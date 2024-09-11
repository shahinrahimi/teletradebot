package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

func (b *Bot) StartBinanceService(ctx context.Context) {
	go b.startUserDataStream724(ctx)
}

func (b *Bot) startUserDataStream(ctx context.Context) {
	futures.UseTestnet = b.bc.UseTestnet
	listenKey, err := b.bc.Client.NewStartUserStreamService().Do(ctx)
	if err != nil {
		b.l.Fatalf("Error starting user data stream : %v", err)
	}
	b.l.Printf("ListenKey: %s", listenKey)
	doneC, stopC, err := futures.WsUserDataServe(listenKey, b.wsHandler, b.errHandler)
	if err != nil {
		log.Fatalf("Error establishing WebSocket connection: %v", err)
	}
	b.l.Println("WebSocket connection established. Awaiting events...")

	// Keep the connection alive by sending a ping every 30 minutes
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				err := b.bc.Client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(ctx)
				if err != nil {
					b.l.Printf("Error keeping user data stream alive: %v", err)
				} else {
					fmt.Println("User data stream kept alive")
				}
			case <-doneC:
				return
			}
		}
	}()

	// Keep the connection alive until manually stopped
	<-doneC
	close(stopC)
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
			log.Printf("Error connecting WebSocket: %v", err)
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
						fmt.Println("User data stream kept alive")
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

func (b *Bot) wsHandler(event *futures.WsUserDataEvent) {
	// handle order trade events
	// fmt.Println("got an event")
	b.handleOrderTradeUpdate(event.OrderTradeUpdate)
}

func (b *Bot) errHandler(err error) {
	b.l.Printf("WebSocket error: %v", err)
}
