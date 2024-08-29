package models

import (
	"fmt"
	"time"
)

type Trade struct {
	ID                int
	OrderID           string // OrderID for placed order from binance api
	State             string
	Account           string
	Side              string
	Pair              string
	Candle            string
	Offset            float32
	SizePercent       int
	SLPercent         int
	TPPercent         int
	ReverseMultiplier int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type KeyTrade struct{}

const (
	CREATE_TABLE_TRADES string = `
		CREATE TABLE IF NOT EXISTS trades (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
			order_id TEXT,
			state TEXT NOT NULL,
			account TEXT NOT NULL,
			pair TEXT NOT NULL,
			side TEXT NOT NULL,
			candle TEXT NOT NULL,
			offset REAL NOT NULL,
			size_percent INTEGER NOT NULL,
			sl_percent INTEGER NOT NULL,
			tp_percent INTEGER NOT NULL,
			reverse_multiplier INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)	
	`
	SELECT_COUNT_TRADES string = `SELECT COUNT(*) FROM trades`
	SELECT_TRADES       string = `SELECT * FROM trades`
	SELECT_TRADE        string = `SELECT * FROM trades WHERE id = ?`
	INSERT_TRADE        string = `INSERT INTO trades (state, account, pair, side, candle, offset, size_percent, sl_percent, tp_percent, reverse_multiplier) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ? ,?)`
	DELETE_TRADE        string = `DELETE FROM trades WHERE id = ?`
	UPDATE_TRADE        string = `UPDATE trades SET state = ? WHERE id = ?`
	UPDATE_ORDER_ID     string = `UPDATE trades SET order_id WHERE id = ?`
)

// ToArgs returns state, account, pair, side, candle, offset, size, stop_percent, target_percent and reverse_multiplier as value
// use for inserting to DB
func (t *Trade) ToArgs() []interface{} {
	return []interface{}{t.State, t.Account, t.Pair, t.Side, t.Candle, t.Offset, t.SizePercent, t.SLPercent, t.TPPercent, t.ReverseMultiplier}
}

// ToFeilds returns id, state, account, pair, side, candle, offset, size, stop_percent, target_percent, reverse_multiplier, created_at and updated_at as reference
// use for scanning from DB
func (t *Trade) ToFelids() []interface{} {
	return []interface{}{&t.ID, &t.State, &t.Account, &t.Pair, &t.Side, &t.Candle, &t.Offset, &t.SizePercent, &t.SLPercent, &t.TPPercent, &t.ReverseMultiplier, &t.CreatedAt, &t.UpdatedAt}
}

func (t *Trade) ToListString() string {
	return fmt.Sprintf("id: %d [%s] %s %s %s %s", t.ID, t.Account, t.Pair, t.Side, t.Candle, t.State)
}

func (t *Trade) ToViewString() string {
	return fmt.Sprintf("id: %d\nAccount: %s\nPair: %s\nSide: %s\nCandle: %s\nOffset: %f\nSizePercent: %d\nSLPercent: %d\nTPPercent: %d\nRM: %d", t.ID, t.Account, t.Pair, t.Side, t.Candle, t.Offset, t.SizePercent, t.SLPercent, t.TPPercent, t.ReverseMultiplier)
}
