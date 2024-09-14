package bot

import "github.com/shahinrahimi/teletradebot/models"

func (b *Bot) findTradeWithAnyOrderID(orderID string) (t *models.Trade, isOrderID bool, isTPOrderID bool, isSLOrderID bool, err error) {
	ts, err := b.s.GetTrades()
	if err != nil {
		return nil, false, false, false, err
	}
	for _, t := range ts {
		switch {
		case t.OrderID == orderID:
			return t, true, false, false, nil
		case t.TPOrderID == orderID:
			return t, false, true, false, nil
		case t.SLOrderID == orderID:
			return t, false, false, true, nil
		default:
			continue
		}
	}
	return nil, false, false, false, nil
}

func (b *Bot) getAllUniqueTrades(account string, side string, state string) (map[string]models.Trade, error) {
	uniqueTrades := make(map[string]models.Trade)
	ts, err := b.s.GetTrades()
	if err != nil {
		return nil, err
	}
	for _, t := range ts {
		if t.Account == account && t.Side == side && t.State == state {
			uniqueTrades[t.Symbol] = *t
		}
	}
	return uniqueTrades, nil
}

func (b *Bot) getAllUniqueRawTrades(rawTrades []models.Trade, account string, side string) (map[string]models.Trade, error) {
	uniqueTrades := make(map[string]models.Trade)
	for _, t := range rawTrades {
		if t.Account == account && t.Side == side {
			uniqueTrades[t.Symbol] = t
		}
	}
	return uniqueTrades, nil
}

func (b *Bot) getAllTrades(account string, side string) ([]models.Trade, error) {
	var trades []models.Trade
	ts, err := b.s.GetTrades()
	if err != nil {
		return nil, err
	}
	for _, t := range ts {
		if t.Account == account && t.Side == side {
			trades = append(trades, *t)
		}
	}

	return trades, nil
}
