package store

import (
	"database/sql"
	"log"
	"os"

	"gihub.com/shahinrahimi/teletradebot/models"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	l  *log.Logger
	db *sql.DB
}

type Storage interface {
	CreateTrade(o *models.Trade) error
	GetTrade(id int) (*models.Trade, error)
	GetTrades() ([]*models.Trade, error)
	DeleteTrade(id int) error
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
