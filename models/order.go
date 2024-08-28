package models

import (
	"fmt"
	"time"
)

type Order struct {
	ID                int
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
			size_percent INTEGER NOT NULL,
			sl_percent INTEGER NOT NULL,
			tp_percent INTEGER NOT NULL,
			reverse_multiplier INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)	
	`
	SELECT_COUNT_ORDERS string = `SELECT COUNT(*) FROM orders`
	SELECT_ORDERS       string = `SELECT * FROM orders`
	SELECT_ORDER        string = `SELECT * FROM orders WHERE id = ?`
	INSERT_ORDER        string = `INSERT INTO orders (state, account, pair, side, candle, offset, size_percent, sl_percent, tp_percent, reverse_multiplier) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ? ,?)`
	DELETE_ORDER        string = `DELETE FROM orders WHERE id = ?`
	UPDATE_ORDER        string = `UPDATE orders SET state = ? WHERE id = ?`
	ACCOUNT_B           string = `Binance`
	ACCOUNT_M           string = `Bitmex`
	SIDE_L              string = `LONG`
	SIDE_S              string = `SHORT`
	CANDLE_1MIN         string = `1m`
	CANDLE_3MIN         string = `3m`
	CANDLE_5MIN         string = `5m`
	CANDLE_15MIN        string = `15m`
	CANDLE_30MIN        string = `30m`
	CANDLE_1H           string = `1h`
	CANDLE_2H           string = `2h`
	CANDLE_4H           string = `4h`
	CANDLE_6H           string = `6h`
	CANDLE_8H           string = `8h`
	CANDLE_12H          string = `12h`
	CANDLE_1D           string = `1d`
	CANDLE_3D           string = `3d`
	CANDLE_1W           string = `1w`
	CANDLE_1M           string = `1M`
	STATE_IDLE          string = `idle`
	STATE_PLACED        string = `placed`
	STATE_FILLED        string = `filled`
	STATE_REVERTING     string = `reverting`
)

// ToArgs returns state, account, pair, side, candle, offset, size, stop_percent, target_percent and reverse_multiplier as value
// use for inserting to DB
func (o *Order) ToArgs() []interface{} {
	return []interface{}{o.State, o.Account, o.Pair, o.Side, o.Candle, o.Offset, o.SizePercent, o.SLPercent, o.TPPercent, o.ReverseMultiplier}
}

// ToFeilds returns id, state, account, pair, side, candle, offset, size, stop_percent, target_percent, reverse_multiplier, created_at and updated_at as reference
// use for scanning from DB
func (o *Order) ToFelids() []interface{} {
	return []interface{}{&o.ID, &o.State, &o.Account, &o.Pair, &o.Side, &o.Candle, &o.Offset, &o.SizePercent, &o.SLPercent, &o.TPPercent, &o.ReverseMultiplier, &o.CreatedAt, &o.UpdatedAt}
}

func (o *Order) ToListString() string {
	return fmt.Sprintf("id: %d [%s] %s %s %s %s", o.ID, o.Account, o.Pair, o.Side, o.Candle, o.State)
}

func (o *Order) ToViewString() string {
	return fmt.Sprintf("id: %d\nAccount: %s\nPair: %s\nSide: %s\nCandle: %s\nOffset: %f\nSizePercent: %d\nSLPercent: %d\nTPPercent: %d\nRM: %d", o.ID, o.Account, o.Pair, o.Side, o.Candle, o.Offset, o.SizePercent, o.SLPercent, o.TPPercent, o.ReverseMultiplier)
}
