package binance

import "context"

func (bc *BinanceClient) GetListenKey(ctx context.Context) (string, error) {
	return bc.Client.NewStartUserStreamService().Do(ctx)
}
