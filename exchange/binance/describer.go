package binance

import (
	"context"
	"time"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) FetchDescriber(ctx context.Context, t *models.Trade) (*models.Describer, error) {
	k, err := bc.getLastClosedKline(ctx, t)
	if err != nil {
		return nil, err
	}
	s, err := bc.getSymbol(t)
	if err != nil {
		return nil, err
	}
	sp, err := bc.calculateStopPrice(t, k, s)
	if err != nil {
		return nil, err
	}
	sl, err := bc.calculateStopLossPrice(t, k, s, sp)
	if err != nil {
		return nil, err
	}
	tp, err := bc.calculateTakeProfitPrice(t, k, s, sp)
	if err != nil {
		return nil, err
	}

	from := utils.ConvertTime(k.OpenTime)
	till := utils.ConvertTime(k.CloseTime).Add(time.Second)

	return &models.Describer{
		From:  from,
		Till:  till,
		Open:  k.Open,
		Close: k.Close,
		High:  k.High,
		Low:   k.Low,
		SP:    sp,
		TP:    tp,
		SL:    sl,
	}, nil
}
