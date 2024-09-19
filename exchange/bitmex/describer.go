package bitmex

import (
	"context"
	"fmt"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (mc *BitmexClient) FetchDescriber(ctx context.Context, t *models.Trade) (*models.Describer, error) {
	timeframeDur, err := types.GetDuration(t.Timeframe)
	if err != nil {
		return nil, err
	}
	k, err := mc.GetLastClosedCandle(t.Symbol, timeframeDur)
	if err != nil {
		return nil, err
	}

	// sp, err := bc.calculateStopPrice(t, k, s)
	// if err != nil {
	// 	return nil, err
	// }
	// sl, err := bc.calculateStopLossPrice(t, k, s, sp)
	// if err != nil {
	// 	return nil, err
	// }
	// tp, err := bc.calculateTakeProfitPrice(t, k, s, sp)
	// if err != nil {
	// 	return nil, err
	// }
	var spFloat float64
	if t.Side == types.SIDE_L {
		spFloat = k.High + t.Offset
	} else {
		spFloat = k.Low - t.Offset
	}

	return &models.Describer{
		From:  k.OpenTime,
		Till:  k.CloseTime,
		Open:  fmt.Sprintf("%.5f", k.Open),
		Close: fmt.Sprintf("%.5f", k.Close),
		High:  fmt.Sprintf("%.5f", k.High),
		Low:   fmt.Sprintf("%.5f", k.Low),
		SP:    fmt.Sprintf("%.5f", spFloat),
		TP:    "0",
		SL:    "0",
	}, nil
}
