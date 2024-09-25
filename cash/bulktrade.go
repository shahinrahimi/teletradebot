package cash

import (
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (c *Cash) GetAllUniqueTrades(account types.ExchangeType, side types.SideType, state types.StateType) []*models.Trade {
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

func (c *Cash) GetAllUniqueRawTrades(rawTrades []models.Trade, account types.ExchangeType, side types.SideType) []*models.Trade {
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

func (c *Cash) GetAllTrades(account types.ExchangeType, side types.SideType) []*models.Trade {
	var trades []*models.Trade
	ts := c.GetTrades()
	for _, t := range ts {
		if t.Account == account && t.Side == side {
			trades = append(trades, t)
		}
	}
	return trades
}
