package models

import (
	"fmt"
	"time"
)

type Trade struct {
	ID                int
	OrderID           string // OrderID for placed order from binance api
	UserID            int64  // id of telegram user
	State             string
	Account           string
	Side              string
	Symbol            string
	Candle            string  // 1h 4h 15m etc
	Offset            float64 // offset for placing the order defines in usdt amount 1 is 1$ 0.1 is 0.1$ etc
	SizePercent       int     // like 1, 2, 3 or 5
	SLPercent         int     // like 101, 105 => calculate base on kline the candle before last (e.g 100 means 100% (range) range = high - low before last candle)
	TPPercent         int     // like 101, 105 => calculate base on kline the candle before last (e.g 104 means 105% (range) range = high - low before last candle)
	ReverseMultiplier int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type KeyTrade struct{}

const (
	CREATE_TABLE_TRADES string = `
		CREATE TABLE IF NOT EXISTS trades (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
			order_id TEXT NOT NULL DEFAULT '',
			user_id INTEGER NOT NULL,
			state TEXT NOT NULL,
			account TEXT NOT NULL,
			symbol TEXT NOT NULL,
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
	SELECT_COUNT_TRADES     string = `SELECT COUNT(*) FROM trades`
	SELECT_TRADES           string = `SELECT * FROM trades`
	SELECT_TRADE            string = `SELECT * FROM trades WHERE id = ?`
	SELECT_TRADE_BY_OrderID string = `SELECT * FROM trades WHERE order_id = ?`
	INSERT_TRADE            string = `INSERT INTO trades (user_id, state, account, symbol, side, candle, offset, size_percent, sl_percent, tp_percent, reverse_multiplier) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ? ,?)`
	DELETE_TRADE            string = `DELETE FROM trades WHERE id = ?`
	UPDATE_TRADE            string = `UPDATE trades SET order_id = ?, state = ?, updated_at = ? WHERE id = ?`
	UPDATE_ORDER_ID         string = `UPDATE trades SET order_id WHERE id = ?`
)

// ToArgs returns state, account, symbol, side, candle, offset, size, stop_percent, target_percent and reverse_multiplier as value
// use for inserting to DB
func (t *Trade) ToArgs() []interface{} {
	return []interface{}{t.UserID, t.State, t.Account, t.Symbol, t.Side, t.Candle, t.Offset, t.SizePercent, t.SLPercent, t.TPPercent, t.ReverseMultiplier}
}

// ToUpdatedArgs returns order_id, state, updated_at and id as value
// use for updating record in DB
func (t *Trade) ToUpdatedArgs() []interface{} {
	return []interface{}{t.OrderID, t.State, t.UpdatedAt, t.ID}
}

// ToFields returns id, order_id, user_id, state, account, symbol, side, candle, offset, size, stop_percent, target_percent, reverse_multiplier, created_at and updated_at as reference
// use for scanning from DB
func (t *Trade) ToFelids() []interface{} {
	return []interface{}{&t.ID, &t.OrderID, &t.UserID, &t.State, &t.Account, &t.Symbol, &t.Side, &t.Candle, &t.Offset, &t.SizePercent, &t.SLPercent, &t.TPPercent, &t.ReverseMultiplier, &t.CreatedAt, &t.UpdatedAt}
}

func (t *Trade) ToListString() string {
	return fmt.Sprintf("id: %d [%s] %s %s %s %s", t.ID, t.Account, t.Symbol, t.Side, t.Candle, t.State)
}

func (t *Trade) ToViewString() string {
	return fmt.Sprintf("id: %d\nAccount: %s\nSymbol: %s\nSide: %s\nCandle: %s\nOffset: %f\nSizePercent: %d\nSLPercent: %d\nTPPercent: %d\nRM: %d", t.ID, t.Account, t.Symbol, t.Side, t.Candle, t.Offset, t.SizePercent, t.SLPercent, t.TPPercent, t.ReverseMultiplier)
}
