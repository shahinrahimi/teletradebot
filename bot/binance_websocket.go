package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

func (b *Bot) StartBinanceService(ctx context.Context) {
	go b.startUserDataStream(ctx)
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

	// Keep the connection alive by sending a ping every 30 miutes
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				err := b.bc.Client.NewKeepaliveUserStreamService().ListenKey(listenKey)
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

func (b *Bot) wsHandler(event *futures.WsUserDataEvent) {
	// handle order trade events
	// fmt.Println("got an event")
	b.handleOrderTradeUpdate(event.OrderTradeUpdate)
}

func (b *Bot) errHandler(err error) {
	b.l.Printf("WebSocket error: %v", err)
}
