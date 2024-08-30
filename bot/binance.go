package bot

import (
	"gihub.com/shahinrahimi/teletradebot/utils"
	"github.com/adshao/go-binance/v2/futures"
)

func (b *Bot) StartBinanceService() error {
	if err := b.bc.UpdateTickers(); err != nil {
		b.l.Printf("error updating tickers for binance : %v", err)
		return err
	}
	b.l.Printf("Total pairs found for binance: %d", len(b.bc.Symbols))

	if err := b.bc.UpdateListenKey(); err != nil {
		b.l.Printf("error updating listenKey for binance : %v", err)
		return err
	}
	b.l.Printf("ListenKey acquired: %s", b.bc.ListenKey)
	return nil
}

func (b *Bot) StartWsBinanceService() {
	go b.startUserDataStream()
}

func (b *Bot) startUserDataStream() {
	futures.UseTestnet = b.bc.UseTestnet

	doneC, _, err := futures.WsUserDataServe(b.bc.ListenKey, b.wsHandler, b.errHandler)
	if err != nil {
		b.l.Printf("error startUserDataStream: %v", err)
		return
	}

	b.l.Println("WebSocket connection established. Listening for events...")
	<-doneC
}

func (b *Bot) handleOrderTradeUpdate(f futures.WsOrderTradeUpdate) {
	utils.PrintStructFields(f)
}

func (b *Bot) wsHandler(event *futures.WsUserDataEvent) {
	// handle order trade events
	b.handleOrderTradeUpdate(event.OrderTradeUpdate)
}

func (b *Bot) errHandler(err error) {
	b.l.Printf("handling ws error: %v", err)
}
