package store

import (
	"database/sql"

	"gihub.com/shahinrahimi/teletradebot/models"
)

func (s *SqliteStore) CreateTrade(t *models.Trade) error {
	if _, err := s.db.Exec(models.INSERT_TRADE, t.ToArgs()...); err != nil {
		s.l.Printf("error creating a new trade: %v", err)
		return err
	}
	return nil
}

func (s *SqliteStore) GetTrade(id int) (*models.Trade, error) {
	var t models.Trade
	if err := s.db.QueryRow(models.SELECT_TRADE, id).Scan(t.ToFelids()...); err != nil {
		if err != sql.ErrNoRows {
			s.l.Printf("error selecting a trade from DB: %v", err)
		}
		return nil, err
	}
	return &t, nil
}

func (s *SqliteStore) GetTradeByOrderID(order_id string) (*models.Trade, error) {
	var t models.Trade
	if err := s.db.QueryRow(models.SELECT_TRADE_BY_OrderID, order_id).Scan(t.ToFelids()...); err != nil {
		if err != sql.ErrNoRows {
			s.l.Printf("trade not found with desired order_id: %s", order_id)
		}
		return nil, err
	}
	return &t, nil
}

func (s *SqliteStore) GetTrades() ([]*models.Trade, error) {
	rows, err := s.db.Query(models.SELECT_TRADES)
	var os []*models.Trade
	for rows.Next() {
		var o models.Trade
		if err := rows.Scan(o.ToFelids()...); err != nil {
			s.l.Printf("error selecting a trade from DB: %v", err)
			continue
		}
		os = append(os, &o)
	}
	return os, err
}

func (s *SqliteStore) DeleteTrade(id int) error {
	if _, err := s.db.Exec(models.DELETE_TRADE, id); err != nil {
		s.l.Printf("error deleting a order from DB: %v", err)
		return err
	}
	return nil
}

func (s *SqliteStore) UpdateTrade(t *models.Trade) error {
	if _, err := s.db.Exec(models.UPDATE_TRADE, t.ToUpdatedArgs()...); err != nil {
		s.l.Printf("error updating trade to DB: %v", err)
		return err
	}
	return nil
}
