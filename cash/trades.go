package cash

import (
	"log"
	"sync"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

var (
	mu sync.RWMutex
)

func (c *Cash) StorageCreateTrade(t *models.Trade) (int64, error) {
	tradeID, err := c.s.CreateTrade(t)
	if err != nil {
		return 0, err
	}
	c.AddTrade(t)
	return tradeID, nil
}

func (c *Cash) StorageRemoveTrade(id int64) error {
	err := c.s.DeleteTrade(id)
	if err != nil {
		return err
	}
	c.RemoveTrade(id)
	return nil
}

func (c *Cash) GetTrade(id int64) *models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	t, exist := c.trades[id]
	if !exist {
		log.Panicf("trade not found: %d", id)
	}
	return t
}

func (c *Cash) AddTrade(t *models.Trade) {
	mu.Lock()
	defer mu.Unlock()
	c.trades[t.ID] = t
}

func (c *Cash) RemoveTrade(id int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(c.trades, id)
}

func (c *Cash) GetTrades() []*models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	var ts []*models.Trade
	for _, t := range c.trades {
		ts = append(ts, t)
	}
	return ts
}

func (c *Cash) UpdateTradeIdle(id int64) {
	t := c.GetTrade(id)
	t.State = types.STATE_IDLE
	t.OrderID = ""
	t.SLOrderID = ""
	t.TPOrderID = ""
	c.trades[id] = t
}

func (c *Cash) UpdateTradePlaced(id int64, orderID string) {
	t := c.GetTrade(id)
	t.State = types.STATE_PLACED
	t.OrderID = orderID
}

func (c *Cash) UpdateTradeTPOrder(id int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t := c.GetTrade(id)
	t.TPOrderID = orderID
}

func (c *Cash) UpdateTradeSLOrder(id int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t := c.GetTrade(id)
	t.SLOrderID = orderID
}

func (c *Cash) UpdateTradeFilled(id int64) {
	mu.Lock()
	defer mu.Unlock()
	t := c.GetTrade(id)
	t.State = types.STATE_FILLED
}

func (c *Cash) UpdateTradeStopped(id int64) {
	mu.Lock()
	defer mu.Unlock()
	t := c.GetTrade(id)
	t.State = types.STATE_STOPPED
}

func (c *Cash) UpdateTradeProfited(id int64) {
	mu.Lock()
	defer mu.Unlock()
	t := c.GetTrade(id)
	t.State = types.STATE_PROFITED
}

func (c *Cash) UpdateTradeCanceled(id int64) {
	mu.Lock()
	defer mu.Unlock()
	t := c.GetTrade(id)
	t.State = types.STATE_CANCELED
}

func (c *Cash) GetTradeByOrderID(orderID string) *models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	for _, t := range c.trades {
		if t.OrderID == orderID {
			return t
		}
	}
	return nil
}

func (c *Cash) GetTradeBySLOrderID(orderID string) *models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	for _, t := range c.trades {
		if t.SLOrderID == orderID {
			return t
		}
	}
	return nil
}

func (c *Cash) GetTradeByTPOrderID(orderID string) *models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	for _, t := range c.trades {
		if t.TPOrderID == orderID {
			return t
		}
	}
	return nil
}

func (c *Cash) GetTradeByAnyOrderID(orderID string) (*models.Trade, types.OrderIDType) {
	mu.RLock()
	defer mu.RUnlock()
	for _, t := range c.trades {
		switch orderID {
		case t.OrderID:
			return t, types.OrderIDTypeMain
		case t.SLOrderID:
			return t, types.OrderIDTypeStopLoss
		case t.TPOrderID:
			return t, types.OrderIDTypeTakeProfit
		default:
			continue
		}
	}
	return nil, types.OrderIDTypeNone
}
