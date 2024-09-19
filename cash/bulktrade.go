package cash

import "github.com/shahinrahimi/teletradebot/models"

func (c *Cash) GetAllUniqueTrades(account string, side string, state string) []*models.Trade {
	trades := make([]*models.Trade, 0)
	uniqueSymbolTrades := make(map[string]models.Trade, 0)
	ts := c.GetTrades()
	for _, t := range ts {
		if t.Account == account && t.Side == side && t.State == state {
			uniqueSymbolTrades[t.Symbol] = *t
		}
	}
	for _, t := range uniqueSymbolTrades {
		trades = append(trades, &t)
	}
	return trades
}

func (c *Cash) GetAllUniqueRawTrades(rawTrades []models.Trade, account string, side string) []*models.Trade {
	trades := make([]*models.Trade, 0)
	uniqueTrades := make(map[string]models.Trade, 0)
	for _, t := range rawTrades {
		if t.Account == account && t.Side == side {
			uniqueTrades[t.Symbol] = t
		}
	}
	for _, t := range uniqueTrades {
		trades = append(trades, &t)
	}
	return trades
}

func (c *Cash) GetAllTrades(account string, side string) []*models.Trade {
	var trades []*models.Trade
	ts := c.GetTrades()
	for _, t := range ts {
		if t.Account == account && t.Side == side {
			trades = append(trades, t)
		}
	}
	return trades
}
