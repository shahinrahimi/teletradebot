package cash

import (
	"log"

	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/store"
)

type Cash struct {
	trades map[int64]*models.Trade
	l      *log.Logger
	s      store.Storage
}

func NewCash(s store.Storage, l *log.Logger) *Cash {
	trades := make(map[int64]*models.Trade, 0)
	ts, err := s.GetTrades()
	if err != nil {
		log.Panic(err)
	}
	for _, t := range ts {
		trades[t.ID] = t
	}
	return &Cash{
		trades: trades,
		l:      l,
		s:      s,
	}
}