package binance

import "context"

func (bc *BinanceClient) SetLeverage(ctx context.Context, symbol string, leverage int) error {
	_, err := bc.client.NewChangeLeverageService().Symbol(symbol).Leverage(leverage).Do(ctx)
	bc.client.NewGetLeverageBracketService()
	if err != nil {
		return err
	}
	return nil
}

func (bc *BinanceClient) ReadAllLeverageSymbols(ctx context.Context) ([]string, error) {
	res, err := bc.client.NewGetLeverageBracketService().Do(ctx)
	if err != nil {
		return nil, err
	}
	var symbols []string
	for _, s := range res {
		symbols = append(symbols, s.Symbol)
		bc.l.Printf("symbol: %s bracket: %d\n", s.Symbol, s.Brackets)
		for _, b := range s.Brackets {
			bc.l.Printf("  bracket: %d InitialLeverage: %d\n", b.Bracket, b.InitialLeverage)
		}
	}
	return symbols, nil
}
