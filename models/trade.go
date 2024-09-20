package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/shahinrahimi/teletradebot/types"
)

type Trade struct {
	ID     int64 // Unique identifier for the trade.
	UserID int64 // ID of the Telegram user associated with the trade.
	ChatID int64 // Chat ID in Telegram where trade created for communication.

	OrderID   string // OrderID for the placed order from the Binance API or Bitmex API.
	SLOrderID string // OrderID for the Stop Loss order.
	TPOrderID string // OrderID for the Take Profit order.
	State     string // Current state of the trade.

	Account           string    // Trading account associated with the trade.
	Side              string    // Side of the trade (e.g., buy or sell).
	Symbol            string    // Trading pair symbol (e.g., BTCUSDT).
	Timeframe         string    // Timeframe of the candle (e.g., 1h, 4h, 15m).
	Offset            float64   // Offset for placing the order, defined in USDT amount (e.g., 1 for $1, 0.1 for $0.1).
	Size              int       // Size of the trade as percentage (e.g., 1, 2, 3, or 5).
	StopLoss          int       // Stop Loss percentage based on the range of the candle before the last (e.g., 100 for 100% of the range).
	TakeProfit        int       // Take Profit percentage based on the range of the candle before the last (e.g., 105 for 105% of the range).
	ReverseMultiplier int       // Multiplier used for reversing the trade.
	CreatedAt         time.Time // Timestamp when the trade was created.
	UpdatedAt         time.Time // Timestamp when the trade was last updated.
}

type KeyTrade struct{}

const (
	CREATE_TABLE_TRADES string = `
		CREATE TABLE IF NOT EXISTS trades (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
			user_id INTEGER NOT NULL,
			chat_id INTEGER NOT NULL,
			order_id TEXT NOT NULL DEFAULT '',
			sl_order_id TEXT NOT NULL DEFAULT '',
			tp_order_id TEXT NOT NULL DEFAULT '',
			state TEXT NOT NULL,
			account TEXT NOT NULL,
			symbol TEXT NOT NULL,
			side TEXT NOT NULL,
			timeframe TEXT NOT NULL,
			offset REAL NOT NULL,
			size INTEGER NOT NULL,
			stop_loss INTEGER NOT NULL,
			take_profit INTEGER NOT NULL,
			reverse_multiplier INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)	
	`
	SELECT_COUNT_TRADES         string = `SELECT COUNT(*) FROM trades`
	SELECT_TRADES               string = `SELECT * FROM trades`
	SELECT_TRADE                string = `SELECT * FROM trades WHERE id = ?`
	SELECT_TRADE_BY_ORDER_ID    string = `SELECT * FROM trades WHERE order_id = ?`
	SELECT_TRADE_BY_SL_ORDER_ID string = `SELECT * FROM trades WHERE sl_order_id = ?`
	SELECT_TRADE_BY_TP_ORDER_ID string = `SELECT * FROM trades WHERE tp_order_id = ?`
	INSERT_TRADE                string = `INSERT INTO trades (user_id, chat_id, state, account, symbol, side, timeframe, offset, size, stop_loss, take_profit, reverse_multiplier) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? ,?)`
	DELETE_TRADE                string = `DELETE FROM trades WHERE id = ?`
	UPDATE_TRADE                string = `UPDATE trades SET order_id = ?, sl_order_id = ?, tp_order_id = ?, state = ?, updated_at = ? WHERE id = ?`
)

// ToArgs returns user_id, chat_id, state, account, symbol, side, timeframe, offset, size, stop_loss, take_profit and reverse_multiplier as value
// use for inserting to DB
func (t *Trade) ToArgs() []interface{} {
	return []interface{}{t.UserID, t.ChatID, t.State, t.Account, t.Symbol, t.Side, t.Timeframe, t.Offset, t.Size, t.StopLoss, t.TakeProfit, t.ReverseMultiplier}
}

// ToUpdatedArgs returns order_id, sl_order_id, tp_order_id, state, updated_at and id as value
// use for updating record in DB
func (t *Trade) ToUpdatedArgs() []interface{} {
	return []interface{}{t.OrderID, t.SLOrderID, t.TPOrderID, t.State, t.UpdatedAt, t.ID}
}

// ToFields returns id, user_id, chat_id, order_id, sl_order_id, tp_order_id, state, account, symbol, side, timeframe, offset, size, stop_loss, take_profit, reverse_multiplier, created_at and updated_at as reference
// use for scanning from DB
func (t *Trade) ToFelids() []interface{} {
	return []interface{}{&t.ID, &t.UserID, &t.ChatID, &t.OrderID, &t.SLOrderID, &t.TPOrderID, &t.State, &t.Account, &t.Symbol, &t.Side, &t.Timeframe, &t.Offset, &t.Size, &t.StopLoss, &t.TakeProfit, &t.ReverseMultiplier, &t.CreatedAt, &t.UpdatedAt}
}

func (t *Trade) ToListString() string {
	return fmt.Sprintf("ID: %d [%s] %s %s %s %s", t.ID, t.Account, t.Symbol, t.Side, t.Timeframe, strings.ToUpper(t.State))
}

func (t *Trade) ToViewString() string {
	return fmt.Sprintf("Trade ID: %d\n\nAccount: %s\nSymbol: %s\nSide: %s\nTimeframe: %s\nOffset: $%0.2f\nSize: %d\nSL: %d\nTP: %d\nRM: %d", t.ID, t.Account, t.Symbol, t.Side, t.Timeframe, t.Offset, t.Size, t.StopLoss, t.TakeProfit, t.ReverseMultiplier)
}

func (t *Trade) CalculateStopPrice(high, low float64) (float64, error) {
	var stopPrice float64
	if t.Side == types.SIDE_L {
		stopPrice = high + t.Offset
	} else {
		stopPrice = low - t.Offset
	}
	if stopPrice <= 0 {
		return 0, fmt.Errorf("price cannot be zero or negative")
	}
	return stopPrice, nil
}

func (t *Trade) CalculateStopLossPrice(high, low, basePrice float64) (float64, error) {
	var stopPrice float64
	r := high - low
	if t.Side == types.SIDE_L {
		stopPrice = basePrice - (r * (float64(t.StopLoss)) / 100)
	} else {
		stopPrice = basePrice + (r * (float64(t.StopLoss)) / 100)
	}
	if stopPrice <= 0 {
		return 0, fmt.Errorf("price cannot be zero or negative")
	}
	return stopPrice, nil
}

func (t *Trade) CalculateTakeProfitPrice(high, low, basePrice float64) (float64, error) {
	var stopPrice float64
	r := high - low
	if t.Side == types.SIDE_L {
		stopPrice = basePrice + (r * (float64(t.TakeProfit)) / 100)
	} else {
		stopPrice = basePrice - (r * (float64(t.TakeProfit)) / 100)
	}
	if stopPrice <= 0 {
		return 0, fmt.Errorf("price cannot be zero or negative")
	}
	return stopPrice, nil
}
