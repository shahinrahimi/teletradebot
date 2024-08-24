package models

import "time"

type Order struct {
	ID                int
	State             string
	Account           string
	Pair              string
	Side              string
	Candle            string
	Offset            float32
	Size              float32
	StopPercent       float32
	TargetPercent     float32
	ReverseMultiplier int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type KeyOrder struct{}

const (
	CREATE_TABLE_ORDERS string = `
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
			state TEXT NOT NULL,
			account TEXT NOT NULL,
			pair TEXT NOT NULL,
			side TEXT NOT NULL,
			candle TEXT NOT NULL,
			offset REAL NOT NULL,
			size REAL NOT NULL,
			stop_percent REAL NOT NULL,
			target_percent REAL NOT NULL,
			reverse_multiplier INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)	
	`
	SELECT_COUNT_ORDERS string = `SELECT COUNT(*) FROM orders`
	SELECT_ORDERS       string = `SELECT * FROM orders`
	SELECT_ORDER        string = `SELECT * FROM orders WHERE id = ?`
	INSERT_ORDER        string = `INSERT INTO orders (state, account, pair, side, candle, offset, size, stop_percent, target_percent, reverse_multiplier) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ? ,?)`
	DELETE_ORDER        string = `DELETE FROM orders WHERE id = ?`
	UPDATE_ORDER        string = `UPDATE orders SET state = ? WHERE id = ?`
)

// ToArgs returns state, account, pair, side, candle, offset, size, stop_percent, target_percent and reverse_multiplier as value
// use for inserting to DB
func (o *Order) ToArgs() []interface{} {
	return []interface{}{o.State, o.Account, o.Pair, o.Side, o.Candle, o.Offset, o.Side, o.StopPercent, o.TargetPercent, o.ReverseMultiplier}
}

// ToFeilds returns id, state, account, pair, side, candle, offset, size, stop_percent, target_percent, reverse_multiplier, created_at and updated_at as reference
// use for scanning from DB
func (o *Order) ToFelids() []interface{} {
	return []interface{}{&o.ID, &o.State, &o.Account, &o.Pair, &o.Side, &o.Candle, &o.Offset, &o.Side, &o.StopPercent, &o.TargetPercent, &o.ReverseMultiplier, &o.CreatedAt, &o.UpdatedAt}
}
