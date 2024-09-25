package cash

import (
	"log"
	"sort"
	"sync"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

var (
	mu sync.RWMutex
)

func (c *Cash) StorageCreateTrade(t *models.Trade) (int64, error) {
	tradeID, err := c.s.CreateTrade(t)
	if err != nil {
		return 0, err
	}
	c.AddTrade(t, tradeID)
	return tradeID, nil
}

func (c *Cash) StorageRemoveTrade(ID int64) error {
	err := c.s.DeleteTrade(ID)
	if err != nil {
		return err
	}
	c.RemoveTrade(ID)
	return nil
}

func (c *Cash) GetTrade(ID int64) *models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	t, exist := c.trades[ID]
	if !exist {
		log.Panicf("trade not found: %d", ID)
	}
	return t
}

func (c *Cash) AddTrade(t *models.Trade, ID int64) {
	mu.Lock()
	defer mu.Unlock()
	t.ID = ID
	c.trades[ID] = t
}

func (c *Cash) RemoveTrade(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(c.trades, ID)
}

func (c *Cash) GetTrades() []*models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	var ts []*models.Trade
	for _, t := range c.trades {
		ts = append(ts, t)
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i].ID < ts[j].ID })
	return ts
}

func (c *Cash) UpdateTradeIdle(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.State = types.STATE_IDLE
	t.OrderID = ""
	t.SLOrderID = ""
	t.TPOrderID = ""
	c.trades[ID] = t
}

func (c *Cash) updateTradeMainOrder(ID int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.State = types.STATE_PLACED
	t.OrderID = orderID
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeMainOrder(ID int64, orderIDorOrderRes interface{}) {
	orderID := utils.ExtractOrderIDStr(orderIDorOrderRes)
	c.updateTradeMainOrder(ID, orderID)
}

func (c *Cash) updateTradeTPOrder(ID int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.TPOrderID = orderID
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeSLOrder(ID int64, orderIDorOrderRes interface{}) {
	orderIDStr := utils.ExtractOrderIDStr(orderIDorOrderRes)
	c.updateTradeSLOrder(ID, orderIDStr)
}

func (c *Cash) updateTradeSLOrder(ID int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.SLOrderID = orderID
	c.trades[ID] = t

}

func (c *Cash) UpdateTradeTPOrder(ID int64, resOrder interface{}) {
	orderID := utils.ExtractOrderIDStr(resOrder)
	c.updateTradeTPOrder(ID, orderID)
}

func (c *Cash) updateTradeReverseOrder(ID int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.ReverseOrderID = orderID
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeReverseOrder(ID int64, resOrder interface{}) {
	orderID := utils.ExtractOrderIDStr(resOrder)
	c.updateTradeReverseOrder(ID, orderID)
}

func (c *Cash) updateTradeReverseSLOrder(ID int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.ReverseSLOrderID = orderID
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeReverseSLOrder(ID int64, resOrder interface{}) {
	orderID := utils.ExtractOrderIDStr(resOrder)
	c.updateTradeReverseSLOrder(ID, orderID)
}

func (c *Cash) updateTradeReverseTPOrder(ID int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.ReverseTPOrderID = orderID
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeReverseTPOrder(ID int64, resOrder interface{}) {
	orderID := utils.ExtractOrderIDStr(resOrder)
	c.updateTradeReverseTPOrder(ID, orderID)
}

func (c *Cash) UpdateTradeFilled(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.State = types.STATE_FILLED
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeStopped(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.State = types.STATE_STOPPED
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeProfited(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.State = types.STATE_PROFITED
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeCanceled(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.State = types.STATE_CANCELED
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeClosed(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.State = types.STATE_CLOSED
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeReverting(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.State = types.STATE_REVERTING
	c.trades[ID] = t
}

func (c *Cash) UpdateTradeExpired(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	t, exist := c.trades[ID]
	if !exist {
		c.l.Panicf("trade not found: %d", ID)
	}
	t.State = types.STATE_EXPIRED
	c.trades[ID] = t
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
		case t.ReverseOrderID:
			return t, types.OrderIDTypeReverseMain
		case t.ReverseSLOrderID:
			return t, types.OrderIDTypeReverseStopLoss
		case t.ReverseTPOrderID:
			return t, types.OrderIDTypeReverseTakeProfit
		default:
			continue
		}
	}
	return nil, types.OrderIDTypeNone
}
