package store

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shahinrahimi/teletradebot/models"
)

type SqliteStore struct {
	l  *log.Logger
	db *sql.DB
}

type Storage interface {
	CreateTrade(o *models.Trade) (int64, error)
	GetTrade(id int64) (*models.Trade, error)
	GetTrades() ([]*models.Trade, error)
	GetTradeByOrderID(orderID string) (*models.Trade, error)
	GetTradeBySLOrderID(SLOrderID string) (*models.Trade, error)
	GetTradeByTPOrderID(TPOrderID string) (*models.Trade, error)
	DeleteTrade(id int64) error
	UpdateTrade(t *models.Trade) error
	UpdateTradeFilled(t *models.Trade) error
	UpdateTradeStopped(t *models.Trade) error
	UpdateTradeProfited(t *models.Trade) error
	UpdateTradeCancelled(t *models.Trade) error
	UpdateTradeSLOrder(t *models.Trade, SLOrderID string) error
	UpdateTradeTPOrder(t *models.Trade, TPOrderID string) error
	UpdateTradePlaced(t *models.Trade, orderID string) error
	UpdateTradeIdle(t *models.Trade) error
}

func NewSqliteStore(l *log.Logger) (*SqliteStore, error) {
	if err := os.MkdirAll("db", 0755); err != nil {
		l.Printf("Unable to create a directory for DB: %v", err)
		return nil, err
	}
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		l.Printf("Unable to connect to DB: %v", err)
		return nil, err
	}
	l.Println("DB Connected!")
	return &SqliteStore{
		l:  l,
		db: db,
	}, nil
}

func (s *SqliteStore) Init() error {
	if _, err := s.db.Exec(models.CREATE_TABLE_TRADES); err != nil {
		s.l.Printf("error creating table for orders: %v", err)
		return err
	}
	return nil
}

func (s *SqliteStore) CloseDB() {
	if err := s.db.Close(); err != nil {
		s.l.Printf("error closing db connection: %v", err)
	}
	s.l.Printf("DB Disconnected")
}
