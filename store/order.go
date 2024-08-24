package store

import (
	"gihub.com/shahinrahimi/teletradebot/models"
)

func (s *SqliteStore) CreateOrder(o *models.Order) error {
	if _, err := s.db.Exec(models.INSERT_ORDER, o.ToArgs()...); err != nil {
		s.l.Printf("error creating a new order: %v", err)
		return err
	}
	return nil
}

func (s *SqliteStore) GetOrder(id int) (*models.Order, error) {
	var o models.Order
	if err := s.db.QueryRow(models.SELECT_ORDER, id).Scan(o.ToFelids()...); err != nil {
		s.l.Printf("error selecting a order from DB: %v", err)
		return nil, err
	}
	return &o, nil
}

func (s *SqliteStore) GetOrders() ([]*models.Order, error) {
	rows, err := s.db.Query(models.SELECT_ORDERS)
	var os []*models.Order
	for rows.Next() {
		var o *models.Order
		if err := rows.Scan(o.ToFelids()...); err != nil {
			s.l.Printf("error selecting a order from DB: %v", err)
			continue
		}
		os = append(os, o)
	}
	return os, err
}

func (s *SqliteStore) DeleteOrder(id int) error {
	if _, err := s.db.Exec(models.DELETE_ORDER, id); err != nil {
		s.l.Printf("error deleting a order from DB: %v", err)
		return err
	}
	return nil
}
