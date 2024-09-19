package cash

import (
	"log"
	"sync"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/store"
	"github.com/shahinrahimi/teletradebot/types"
)

var (
	trades = make(map[int64]*models.Trade, 0)
	mu     sync.RWMutex
)

func InitTrades(s store.Storage) {
	ts, err := s.GetTrades()
	if err != nil {
		log.Panic(err)
	}
	for _, t := range ts {
		t.State = types.STATE_IDLE
		t.OrderID = ""
		t.SLOrderID = ""
		t.TPOrderID = ""
		trades[t.ID] = t
	}
}

func GetTrade(id int64) *models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	t, exist := trades[id]
	if !exist {
		log.Panicf("trade not found: %d", id)
	}
	return t
}

func AddTrade(t *models.Trade) {
	mu.Lock()
	defer mu.Unlock()
	trades[t.ID] = t
}

func RemoveTrade(id int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(trades, id)
}

func GetTrades() []*models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	var ts []*models.Trade
	for _, t := range trades {
		ts = append(ts, t)
	}
	return ts
}

func UpdateTradeIdle(id int64) {
	t := GetTrade(id)
	t.State = types.STATE_IDLE
	t.OrderID = ""
	t.SLOrderID = ""
	t.TPOrderID = ""
	trades[id] = t
}

func UpdateTradePlaced(id int64, orderID string) {
	t := GetTrade(id)
	t.State = types.STATE_PLACED
	t.OrderID = orderID
}

func UpdateTradeTPOrder(id int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t := GetTrade(id)
	t.TPOrderID = orderID
}

func UpdateTradeSLOrder(id int64, orderID string) {
	mu.Lock()
	defer mu.Unlock()
	t := GetTrade(id)
	t.SLOrderID = orderID
}

func UpdateTradeFilled(id int64) {
	mu.Lock()
	defer mu.Unlock()
	t := GetTrade(id)
	t.State = types.STATE_FILLED
}

func UpdateTradeStopped(id int64) {
	mu.Lock()
	defer mu.Unlock()
	t := GetTrade(id)
	t.State = types.STATE_STOPPED
}

func UpdateTradeProfited(id int64) {
	mu.Lock()
	defer mu.Unlock()
	t := GetTrade(id)
	t.State = types.STATE_PROFITED
}

func UpdateTradeCanceled(id int64) {
	mu.Lock()
	defer mu.Unlock()
	t := GetTrade(id)
	t.State = types.STATE_CANCELED
}

func GetTradeByOrderID(orderID string) *models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	for _, t := range trades {
		if t.OrderID == orderID {
			return t
		}
	}
	return nil
}

func GetTradeBySLOrderID(orderID string) *models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	for _, t := range trades {
		if t.SLOrderID == orderID {
			return t
		}
	}
	return nil
}

func GetTradeByTPOrderID(orderID string) *models.Trade {
	mu.RLock()
	defer mu.RUnlock()
	for _, t := range trades {
		if t.TPOrderID == orderID {
			return t
		}
	}
	return nil
}

func GetTradeByAnyOrderID(orderID string) (*models.Trade, types.OrderIDType) {
	mu.RLock()
	defer mu.RUnlock()
	for _, t := range trades {
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
