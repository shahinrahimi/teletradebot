package binance

import "context"

func (bc *BinanceClient) GetListenKey(ctx context.Context) (string, error) {
	return bc.client.NewStartUserStreamService().Do(ctx)
}
