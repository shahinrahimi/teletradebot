package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shahinrahimi/teletradebot/types"
)

type Trade struct {
	ID     int64 // Unique identifier for the trade.
	UserID int64 // ID of the Telegram user associated with the trade.
	ChatID int64 // Chat ID in Telegram where trade created for communication.

	OrderID          string // OrderID for the placed order from the Binance API or Bitmex API.
	SLOrderID        string // OrderID for the Stop Loss order.
	TPOrderID        string // OrderID for the Take Profit order.
	ReverseOrderID   string // OrderID for the Reverse order.
	ReverseTPOrderID string // OrderID for the Reverse Take Profit order.
	ReverseSLOrderID string // OrderID for the Reverse Stop Loss order.

	State string // Current state of the trade.

	Account           string        // Trading account associated with the trade.
	Side              string        // Side of the trade (e.g., buy or sell).
	Symbol            string        // Trading pair symbol (e.g., BTCUSDT).
	Timeframe         TimeframeType // Timeframe of the candle (e.g., 1h, 4h, 15m).
	Offset            float64       // Offset for placing the order, defined in USDT amount (e.g., 1 for $1, 0.1 for $0.1).
	Size              int           // Size of the trade as percentage (e.g., 1, 2, 3, or 5).
	StopLossSize      int           // Stop Loss percentage based on the range of the candle before the last (e.g., 100 for 100% of the range).
	TakeProfitSize    int           // Take Profit percentage based on the range of the candle before the last (e.g., 105 for 105% of the range).
	ReverseMultiplier int           // Multiplier used for reversing the trade.
	CreatedAt         time.Time     // Timestamp when the trade was created.
	UpdatedAt         time.Time     // Timestamp when the trade was last updated.
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
			stop_loss_size INTEGER NOT NULL,
			take_profit_size INTEGER NOT NULL,
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
	return []interface{}{t.UserID, t.ChatID, t.State, t.Account, t.Symbol, t.Side, t.Timeframe, t.Offset, t.Size, t.StopLossSize, t.TakeProfitSize, t.ReverseMultiplier}
}

// ToUpdatedArgs returns order_id, sl_order_id, tp_order_id, state, updated_at and id as value
// use for updating record in DB
func (t *Trade) ToUpdatedArgs() []interface{} {
	return []interface{}{t.OrderID, t.SLOrderID, t.TPOrderID, t.State, t.UpdatedAt, t.ID}
}

// ToFields returns id, user_id, chat_id, order_id, sl_order_id, tp_order_id, state, account, symbol, side, timeframe, offset, size, stop_loss, take_profit, reverse_multiplier, created_at and updated_at as reference
// use for scanning from DB
func (t *Trade) ToFelids() []interface{} {
	return []interface{}{&t.ID, &t.UserID, &t.ChatID, &t.OrderID, &t.SLOrderID, &t.TPOrderID, &t.State, &t.Account, &t.Symbol, &t.Side, &t.Timeframe, &t.Offset, &t.Size, &t.StopLossSize, &t.TakeProfitSize, &t.ReverseMultiplier, &t.CreatedAt, &t.UpdatedAt}
}

func (t *Trade) ToListString() string {
	return fmt.Sprintf("ID: %d [%s] %s %s %s %s", t.ID, t.Account, t.Symbol, t.Side, t.Timeframe, strings.ToUpper(t.State))
}

func (t *Trade) ToViewString() string {
	return fmt.Sprintf("Trade ID: %d\n\nAccount: %s\nSymbol: %s\nSide: %s\nTimeframe: %s\nOffset: $%0.2f\nSize: %d\nSL: %d\nTP: %d\nRM: %d", t.ID, t.Account, t.Symbol, t.Side, t.Timeframe, t.Offset, t.Size, t.StopLossSize, t.TakeProfitSize, t.ReverseMultiplier)
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

func (t *Trade) CalculateStopLossPrice(high, low, basePrice float64, reverse bool) (float64, error) {
	var stopPrice float64
	r := high - low
	rAmount := (r * (float64(t.StopLossSize)) / 100)
	if t.Side == types.SIDE_L {
		if !reverse {
			stopPrice = basePrice - rAmount
		} else {
			stopPrice = basePrice + rAmount
		}
	} else {
		if !reverse {
			stopPrice = basePrice + rAmount
		} else {
			stopPrice = basePrice - rAmount
		}
	}
	if stopPrice <= 0 {
		return 0, fmt.Errorf("price cannot be zero or negative")
	}
	return stopPrice, nil
}

func (t *Trade) CalculateTakeProfitPrice(high, low, basePrice float64, reverse bool) (float64, error) {
	var stopPrice float64
	r := high - low
	rAmount := (r * (float64(t.TakeProfitSize)) / 100)
	if t.Side == types.SIDE_L {
		if !reverse {
			stopPrice = basePrice + rAmount
		} else {
			stopPrice = basePrice - rAmount
		}
	} else {
		if !reverse {
			stopPrice = basePrice - rAmount
		} else {
			stopPrice = basePrice + rAmount
		}
	}
	return stopPrice, nil
}

func ParseTrade(tradeArgs []string) (*Trade, error) {
	var t Trade
	if len(tradeArgs) < 9 {
		return nil, fmt.Errorf("insufficient arguments provided; please ensure you have 9 parameters")
	}
	// account it should be string
	// m for bitmex
	// b for binance
	part1 := strings.TrimSpace(strings.ToLower(tradeArgs[0]))
	if len(part1) > 1 || (part1 != "m" && part1 != "b") {
		return nil, fmt.Errorf("invalid account value; use 'm' for BitMEX or 'b' for Binance")
	} else if part1 == "m" {
		t.Account = types.ACCOUNT_M
	} else if part1 == "b" {
		t.Account = types.ACCOUNT_B
	} else {
		// should never happen
		return nil, fmt.Errorf("unexpected internal error")
	}
	// pair
	part2 := strings.TrimSpace(strings.ToUpper(tradeArgs[1]))
	t.Symbol = part2
	// side
	part3 := strings.TrimSpace(strings.ToUpper(tradeArgs[2]))
	if part3 != types.SIDE_L && part3 != types.SIDE_S {
		return nil, fmt.Errorf("invalid side value; please enter 'long' or 'short'")
	} else {
		t.Side = part3
	}
	// candle
	part4 := strings.TrimSpace(tradeArgs[3])
	if !IsValidTimeframe(part4) {
		return nil, fmt.Errorf("invalid timeframe; valid values are: %s", GetValidTimeframesString())
	} else {
		t.Timeframe = TimeframeType(part4)
	}
	// offset
	part5 := strings.TrimSpace(tradeArgs[4])
	offset, err := strconv.ParseFloat(part5, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid offset_entry; please provide a numeric value")
	} else {
		t.Offset = offset
	}
	// size percent
	part6 := strings.TrimSpace(tradeArgs[5])
	size_percent, err := strconv.Atoi(part6)
	if err != nil {
		return nil, fmt.Errorf("invalid size; please provide a percentage value (e.g., 5)")
	} else if size_percent <= 0 || size_percent > 50 {
		return nil, fmt.Errorf("invalid size; please provide a value between 1 and 50")
	} else {
		t.Size = size_percent
	}

	// stop-loss percent
	part7 := strings.TrimSpace(tradeArgs[6])
	stop_percent, err := strconv.Atoi(part7)
	if err != nil {
		return nil, fmt.Errorf("invalid stop-loss percent; please provide a numeric value (e.g., 105)")
	} else if stop_percent < 10 {
		return nil, fmt.Errorf("invalid stop-loss percent; must be 10 or greater")
	} else {
		t.StopLossSize = stop_percent
	}

	// target-point percent
	part8 := strings.TrimSpace(tradeArgs[7])
	target_percent, err := strconv.Atoi(part8)
	if err != nil {
		return nil, fmt.Errorf("invalid target-point percent; please provide a numeric value (e.g., 105)")
	} else if target_percent < 10 {
		return nil, fmt.Errorf("invalid target-point percent; must be 10 or greater")
	} else {
		t.TakeProfitSize = target_percent
	}

	// reverse-multiplier
	part9 := strings.TrimSpace(tradeArgs[8])
	reverse_multiplier, err := strconv.Atoi(part9)
	if err != nil {
		return nil, fmt.Errorf("invalid reverse_multiplier; please provide a value of 1 or 2")
	} else if reverse_multiplier <= 0 || reverse_multiplier > 2 {
		return nil, fmt.Errorf("invalid reverse_multiplier; must be 1 or 2")
	} else {
		t.ReverseMultiplier = reverse_multiplier
	}

	t.State = types.STATE_IDLE

	return &t, nil
}
