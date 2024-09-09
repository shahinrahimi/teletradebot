package store

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(models.CREATE_TABLE_TRADES)
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func createTestTrade() *models.Trade {
	return &models.Trade{
		UserID:            123,
		ChatID:            456,
		State:             types.STATE_IDLE,
		Account:           "test_account",
		Symbol:            "BTCUSDT",
		Side:              "buy",
		Timeframe:         "1h",
		Offset:            1.0,
		Size:              10,
		StopLoss:          5,
		TakeProfit:        10,
		ReverseMultiplier: 2,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
}

func TestCreateTrade(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	// Verify the trade was inserted correctly
	insertedTrade, err := store.GetTrade(trade.ID)
	require.NoError(t, err)
	assert.Equal(t, trade.UserID, insertedTrade.UserID)
	assert.Equal(t, trade.ChatID, insertedTrade.ChatID)
	assert.Equal(t, trade.State, insertedTrade.State)
}

func TestGetTrade(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	retrievedTrade, err := store.GetTrade(1)
	require.NoError(t, err)
	assert.NotNil(t, retrievedTrade)
	assert.Equal(t, trade.ID, retrievedTrade.ID)
}

func TestGetTradeByOrderID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1

	err := store.CreateTrade(trade)
	require.NoError(t, err)

	trade.OrderID = "order123"
	err = store.UpdateTrade(trade)
	require.NoError(t, err)

	retrievedTrade, err := store.GetTradeByOrderID(trade.OrderID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedTrade)
	assert.Equal(t, trade.OrderID, retrievedTrade.OrderID)
}

func TestGetTradeBySLOrderID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1

	err := store.CreateTrade(trade)
	require.NoError(t, err)

	trade.SLOrderID = "order123"
	err = store.UpdateTrade(trade)
	require.NoError(t, err)

	retrievedTrade, err := store.GetTradeBySLOrderID(trade.SLOrderID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedTrade)
	assert.Equal(t, trade.SLOrderID, retrievedTrade.SLOrderID)
}

func TestGetTradeByTPOrderID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1

	err := store.CreateTrade(trade)
	require.NoError(t, err)

	trade.TPOrderID = "order123"
	err = store.UpdateTrade(trade)
	require.NoError(t, err)

	retrievedTrade, err := store.GetTradeByTPOrderID(trade.TPOrderID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedTrade)
	assert.Equal(t, trade.TPOrderID, retrievedTrade.TPOrderID)
}

func TestGetTrades(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade1 := createTestTrade()
	trade1.ID = 1
	trade2 := createTestTrade()
	trade2.ID = 2
	trade2.Symbol = "ETHUSDT"

	err := store.CreateTrade(trade1)
	require.NoError(t, err)

	err = store.CreateTrade(trade2)
	require.NoError(t, err)

	trades, err := store.GetTrades()
	require.NoError(t, err)
	assert.Equal(t, 2, len(trades))
}

func TestDeleteTrade(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	err = store.DeleteTrade(trade.ID)
	require.NoError(t, err)

	_, err = store.GetTrade(trade.ID)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestUpdateTrade(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	trade.OrderID = "updated_order_id"
	trade.State = types.STATE_FILLED
	err = store.UpdateTrade(trade)
	require.NoError(t, err)

	updatedTrade, err := store.GetTrade(trade.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated_order_id", updatedTrade.OrderID)
	assert.Equal(t, types.STATE_FILLED, updatedTrade.State)
}

func TestUpdateTradeFilled(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	err = store.UpdateTradeFilled(trade)
	require.NoError(t, err)

	updatedTrade, err := store.GetTrade(trade.ID)
	require.NoError(t, err)
	assert.Equal(t, types.STATE_FILLED, updatedTrade.State)
}

func TestUpdateTradeSLandTP(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	err = store.UpdateTradeSLandTP(trade, "sl_order_123", "tp_order_456")
	require.NoError(t, err)

	updatedTrade, err := store.GetTrade(trade.ID)
	require.NoError(t, err)
	assert.Equal(t, "sl_order_123", updatedTrade.SLOrderID)
	assert.Equal(t, "tp_order_456", updatedTrade.TPOrderID)
}

func TestUpdateTradePlaced(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	err = store.UpdateTradePlaced(trade, "order123")
	require.NoError(t, err)

	updatedTrade, err := store.GetTrade(trade.ID)
	require.NoError(t, err)
	assert.Equal(t, "order123", updatedTrade.OrderID)
}

func TestUpdateTradeIdle(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	err = store.UpdateTradeIdle(trade)
	require.NoError(t, err)

	updatedTrade, err := store.GetTrade(trade.ID)
	require.NoError(t, err)
	assert.Equal(t, types.STATE_IDLE, updatedTrade.State)
	assert.Empty(t, updatedTrade.OrderID)
	assert.Empty(t, updatedTrade.SLOrderID)
	assert.Empty(t, updatedTrade.TPOrderID)
}

func TestUpdateTradeStopped(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	err = store.UpdateTradeStopped(trade)
	require.NoError(t, err)

	updatedTrade, err := store.GetTrade(trade.ID)
	require.NoError(t, err)
	assert.Equal(t, types.STATE_STOPPED, updatedTrade.State)
}

func TestUpdateTradeProfited(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := SqliteStore{db: db}

	trade := createTestTrade()
	trade.ID = 1
	err := store.CreateTrade(trade)
	require.NoError(t, err)

	err = store.UpdateTradeProfited(trade)
	require.NoError(t, err)

	updatedTrade, err := store.GetTrade(trade.ID)
	require.NoError(t, err)
	assert.Equal(t, types.STATE_PROFITED, updatedTrade.State)
}
