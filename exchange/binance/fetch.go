package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) fetchBalances(ctx context.Context) ([]*futures.Balance, error) {
	res, err := bc.client.NewGetBalanceService().Do(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
	// for _, balance := range res {
	// 	if balance.Asset == "USDT" {
	// 		return strconv.ParseFloat(balance.Balance, 64)
	// 	}
	// }
	// return 0, fmt.Errorf("asset USDT not found")
}

func (bc *BinanceClient) fetchPrice(ctx context.Context, symbol string) (float64, error) {
	res, err := bc.client.NewListPricesService().Do(ctx)
	if err != nil {
		return 0, err
	}
	for _, sp := range res {
		if sp.Symbol == symbol {
			return strconv.ParseFloat(sp.Price, 64)
		}
	}
	return 0, fmt.Errorf("symbol %s not found", symbol)
}

func (bc *BinanceClient) fetchCollateral(ctx context.Context) (float64, error) {

	var cc = 4
	errChan := make(chan error, cc)
	balancesChan := make(chan []*futures.Balance, 1)
	priceChanBTC := make(chan float64, 1)
	priceChanETH := make(chan float64, 1)
	priceChanBNB := make(chan float64, 1)

	go func() {
		bs, err := bc.fetchBalances(ctx)
		if err != nil {
			errChan <- fmt.Errorf("error fetching balances: %s", err)
		}
		balancesChan <- bs
	}()

	go func() {
		price, err := bc.fetchPrice(ctx, fmt.Sprintf("%sUSDT", "BTC"))
		if err != nil {
			errChan <- fmt.Errorf("error fetching price for BTC: %s", err)
		}
		priceChanBTC <- price
	}()

	go func() {
		price, err := bc.fetchPrice(ctx, fmt.Sprintf("%sUSDT", "ETH"))
		if err != nil {
			errChan <- fmt.Errorf("error fetching price for ETH: %s", err)
		}
		priceChanETH <- price
	}()

	go func() {
		price, err := bc.fetchPrice(ctx, fmt.Sprintf("%sUSDT", "BNB"))
		if err != nil {
			errChan <- fmt.Errorf("error fetching price for BNB: %s", err)
		}
		priceChanBNB <- price
	}()

	var balances []*futures.Balance
	var priceBTC float64
	var priceETH float64
	var priceBNB float64

	for i := 0; i < cc; i++ {
		select {
		case err := <-errChan:
			return 0, err
		case balances = <-balancesChan:
		case priceBTC = <-priceChanBTC:
		case priceETH = <-priceChanETH:
		case priceBNB = <-priceChanBNB:
		}
	}

	var collateral float64
	for _, balance := range balances {
		switch balance.Asset {
		case "USDT":
			parsedFloat, err := strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return 0, fmt.Errorf("error parsing USDT balance: %s", err)
			}
			collateral = collateral + parsedFloat
		case "BTC":
			parsedFloat, err := strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return 0, fmt.Errorf("error parsing BTC balance: %s", err)
			}
			collateral = collateral + (parsedFloat * priceBTC)
		case "ETH":
			parsedFloat, err := strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return 0, fmt.Errorf("error parsing ETH balance: %s", err)
			}
			collateral = collateral + (parsedFloat * priceETH)
		case "BNB":
			parsedFloat, err := strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return 0, fmt.Errorf("error parsing BNB balance: %s", err)
			}
			collateral = collateral + (parsedFloat * priceBNB)
		default:
			bc.DbgChan <- fmt.Sprintf("skipping asset: %s", balance.Asset)
			continue
		}
	}
	bc.DbgChan <- fmt.Sprintf("Collateral calculated: %f", collateral)
	return collateral, nil

}

func (bc *BinanceClient) fetchLastCompletedCandle(ctx context.Context, symbol string, t models.TimeframeType) (*futures.Kline, error) {
	klines, err := bc.client.NewMarkPriceKlinesService().
		Limit(100).
		Interval(string(t)).
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	// loop through klines and return the most recent completely closed candle
	for i := len(klines) - 1; i >= 0; i-- {
		candleCloseTime := utils.ConvertTime(klines[i].CloseTime)
		// check if close time in the past
		if (time.Until(candleCloseTime)) < 0 {
			return klines[i], nil
		}
	}

	return nil, fmt.Errorf("failed to locate before last candle")
}

func (bc *BinanceClient) fetchSymbolBracket(ctx context.Context, symbol string) ([]futures.Bracket, error) {
	res, err := bc.client.NewGetLeverageBracketService().Symbol(symbol).Do(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range res {
		if s.Symbol == symbol {
			return s.Brackets, nil
		}
	}
	return nil, fmt.Errorf("symbol %s not found", symbol)
}

func (bc *BinanceClient) FetchInterpreter(ctx context.Context, t *models.Trade) (*models.Interpreter, error) {
	bc.DbgChan <- fmt.Sprintf("Fetching interpreter for trade: %d", t.ID)
	var cc = 7
	errChan := make(chan error, cc)
	balancesChan := make(chan []*futures.Balance, 1)
	multiAssetChan := make(chan bool, 1)
	collateralChan := make(chan float64, 1)
	priceChan := make(chan float64, 1)
	candleChan := make(chan *futures.Kline, 1)
	symbolChan := make(chan *futures.Symbol, 1)
	symbolBracketChan := make(chan []futures.Bracket, 1) // contain leverage info

	go func() {
		candle, err := bc.fetchLastCompletedCandle(ctx, t.Symbol, t.Timeframe)
		if err != nil {
			errChan <- err
			candleChan <- nil
			return
		}
		candleChan <- candle
		bc.DbgChan <- fmt.Sprintf("Using candle: %v", candle)
	}()
	go func() {
		symbol, err := bc.GetSymbol(t.Symbol)
		if err != nil {
			errChan <- err
		}
		symbolChan <- symbol
		bc.DbgChan <- fmt.Sprintf("Using symbol: %v", symbol)
	}()
	go func() {
		balances, err := bc.fetchBalances(ctx)
		if err != nil {
			errChan <- err
		}
		balancesChan <- balances
		bc.DbgChan <- fmt.Sprintf("Using balance: %v", balances)
	}()
	go func() {
		price, err := bc.fetchPrice(ctx, t.Symbol)
		if err != nil {
			errChan <- err
		}
		priceChan <- price
		bc.DbgChan <- fmt.Sprintf("Using price: %f", price)
	}()

	go func() {
		symbolBracket, err := bc.fetchSymbolBracket(ctx, t.Symbol)
		if err != nil {
			errChan <- err
			symbolBracketChan <- nil
			return
		}
		symbolBracketChan <- symbolBracket
		bc.DbgChan <- fmt.Sprintf("Using symbol bracket: %v", symbolBracket)
	}()

	go func() {
		isMultiAsset, err := bc.CheckMultiAssetMode(ctx)
		if err != nil {
			errChan <- err
		}
		multiAssetChan <- isMultiAsset
	}()

	go func() {
		collateral, err := bc.fetchCollateral(ctx)
		if err != nil {
			errChan <- err
		}
		collateralChan <- collateral
		bc.DbgChan <- fmt.Sprintf("Using collateral: %f", collateral)
	}()

	var candle *futures.Kline
	var symbol *futures.Symbol
	var symbolBracket []futures.Bracket
	var balances []*futures.Balance
	var collateralBalance float64
	var isMultiAsset bool
	var price float64
	var err error

	for i := 0; i < cc; i++ {
		select {
		case err = <-errChan:
			return nil, err
		case candle = <-candleChan:
		case balances = <-balancesChan:
		case price = <-priceChan:
		case symbol = <-symbolChan:
		case collateralBalance = <-collateralChan:
		case symbolBracket = <-symbolBracketChan:
		case isMultiAsset = <-multiAssetChan:
		}
	}
	bc.DbgChan <- fmt.Sprintf("bracket length: %d", len(symbolBracket))

	var usingBalance float64
	if !isMultiAsset {

		var b *futures.Balance
		for _, balance := range balances {
			if balance.Asset == symbol.QuoteAsset {
				b = balance
				break
			}
		}
		if b == nil {
			return nil, &types.BotError{
				Msg: fmt.Sprintf("cannot find balance for %s", symbol.QuoteAsset),
			}
		}
		usingBalance, err = strconv.ParseFloat(b.Balance, 64)
		if err != nil {
			return nil, err
		}
		bc.DbgChan <- fmt.Sprintf("Using balance in single asset: %f", usingBalance)
	} else {
		usingBalance = collateralBalance
		bc.DbgChan <- fmt.Sprintf("Using balance in multi asset (collateral): %f", collateralBalance)
	}

	size := usingBalance * float64(t.Size) / 100
	quantity := size / price
	rQuantity := quantity * float64(t.ReverseMultiplier)

	// check quantity
	if quantity == 0 {
		return nil, &types.BotError{
			Msg: fmt.Sprintf("insufficient balance. balance: %f.", usingBalance),
		}
	}

	timeframeDur, err := models.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}
	high, err := strconv.ParseFloat(candle.High, 64)
	if err != nil {
		return nil, err
	}
	low, err := strconv.ParseFloat(candle.Low, 64)
	if err != nil {
		return nil, err
	}
	open, err := strconv.ParseFloat(candle.Open, 64)
	if err != nil {
		return nil, err
	}
	close, err := strconv.ParseFloat(candle.Close, 64)
	if err != nil {
		return nil, err
	}
	ep, err := t.CalculateEntryPrice(high, low)
	if err != nil {
		return nil, err
	}
	sl, err := t.CalculateStopLossPrice(high, low, ep, false)
	if err != nil {
		return nil, err
	}
	tp, err := t.CalculateTakeProfitPrice(high, low, ep, false)
	if err != nil {
		return nil, err
	}
	rsl, err := t.CalculateStopLossPrice(high, low, sl, true)
	if err != nil {
		return nil, err
	}
	rtp, err := t.CalculateTakeProfitPrice(high, low, sl, true)
	if err != nil {
		return nil, err
	}

	bc.DbgChan <- fmt.Sprintf("the size is %f balance is %f price is %f quantity is %f reverse quantity is %f", size, usingBalance, price, quantity, rQuantity)

	return &models.Interpreter{
		Balance:         usingBalance,
		Price:           price,
		Quantity:        quantity,
		ReverseQuantity: rQuantity,
		Exchange:        types.ExchangeBinance,

		TradeID:           t.ID,
		Symbol:            t.Symbol,
		Size:              t.Size,
		Side:              t.Side,
		TakeProfitSize:    t.TakeProfitSize,
		StopLossSize:      t.StopLossSize,
		ReverseMultiplier: t.ReverseMultiplier,
		TimeFrameDur:      timeframeDur,

		OpenTime:               utils.ConvertTime(candle.OpenTime),
		CloseTime:              utils.ConvertTime(candle.CloseTime).Add(time.Second),
		Open:                   open,
		Close:                  close,
		High:                   high,
		Low:                    low,
		EntryPrice:             ep,
		StopLossPrice:          sl,
		TakeProfitPrice:        tp,
		ReverseEntryPrice:      sl,
		ReverseStopLossPrice:   rsl,
		ReverseTakeProfitPrice: rtp,
		PricePrecision:         symbol.PricePrecision,
		QuantityPrecision:      symbol.QuantityPrecision,
	}, nil
}
