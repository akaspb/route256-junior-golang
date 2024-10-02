package postgres

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

var (
	ErrorOrderNotFound  = errors.New("order not found")
	ErrorStatusNotFound = errors.New("status not found")
	ErrorPackNotFound   = errors.New("pack not found")
)

type PgStorage struct {
	txManager TransactionManager
}

func NewPgStorage(txManager TransactionManager) *PgStorage {
	return &PgStorage{
		txManager: txManager,
	}
}

func (s *PgStorage) CreateOrder(ctx context.Context, order Order) error {
	tx := s.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
		INSERT INTO orders(id, customer_id, expiry, weight, cost)
			VALUES ($1, $2, $3, $4, $5)
	`, order.ID, order.CustomerID, order.Expiry, order.Weight, order.Cost)

	return err
}

func (s *PgStorage) CreateStatus(ctx context.Context, status Status) error {
	tx := s.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `
		INSERT INTO statuses(order_id, "value", "time") VALUES ($1, $2, $3)
	`, status.OrderID, status.Value, status.Time)

	return err
}

func (s *PgStorage) CreatePack(ctx context.Context, pack Pack) error {
	tx := s.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
		INSERT INTO packs(order_id, name, cost, max_order_weight) VALUES ($1, $2, $3, $4)
	`, pack.OrderID, pack.Name, pack.Cost, pack.MaxOrderWeight)

	return err
}

func (s *PgStorage) GetOrderIDsWhereStatus(ctx context.Context, statusVal models.StatusVal) ([]models.IDType, error) {
	var orderIDs []models.IDType

	tx := s.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orderIDs, `
		SELECT O.id FROM orders O JOIN statuses S ON O.id = S.order_id WHERE S."value" = $1 
	`, statusVal)

	return orderIDs, err
}

func (s *PgStorage) GetOrder(ctx context.Context, orderID models.IDType) (Order, error) {
	var order Order

	tx := s.txManager.GetQueryEngine(ctx)
	err := pgxscan.Get(ctx, tx, &order, `
		SELECT * FROM orders WHERE id = $1
	`, orderID)
	if pgxscan.NotFound(err) {
		return Order{}, ErrorOrderNotFound
	}

	return order, err
}

func (s *PgStorage) GetStatus(ctx context.Context, orderID models.IDType) (Status, error) {
	var status Status

	tx := s.txManager.GetQueryEngine(ctx)
	err := pgxscan.Get(ctx, tx, &status, `
		SELECT * FROM statuses WHERE order_id = $1
	`, orderID)
	if pgxscan.NotFound(err) {
		return Status{}, ErrorStatusNotFound
	}

	return status, err
}

func (s *PgStorage) GetPack(ctx context.Context, orderID models.IDType) (Pack, error) {
	var pack Pack

	tx := s.txManager.GetQueryEngine(ctx)
	err := pgxscan.Get(ctx, tx, &pack, `
		SELECT * FROM packs WHERE order_id = $1
	`, orderID)
	if pgxscan.NotFound(err) {
		return Pack{}, ErrorPackNotFound
	}

	return pack, err
}

func (s *PgStorage) SetStatus(ctx context.Context, status Status) error {
	tx := s.txManager.GetQueryEngine(ctx)

	result, err := tx.Exec(ctx, `
		UPDATE statuses SET "value" = $1, "time" = $2 WHERE order_id = $3
	`, status.Value, status.Time, status.OrderID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrorStatusNotFound
	}

	return nil
}

func (s *PgStorage) DeleteOrder(ctx context.Context, orderId models.IDType) error {
	tx := s.txManager.GetQueryEngine(ctx)

	result, err := tx.Exec(ctx, `
		DELETE FROM orders WHERE id = $1
	`, orderId)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrorOrderNotFound
	}

	return nil
}

func (s *PgStorage) DeletePack(ctx context.Context, orderId models.IDType) error {
	tx := s.txManager.GetQueryEngine(ctx)

	result, err := tx.Exec(ctx, `
		DELETE FROM packs WHERE order_id = $1
	`, orderId)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrorPackNotFound
	}

	return nil
}

func (s *PgStorage) GetCustomerOrdersWithStatus(ctx context.Context, customerId models.IDType, statusVal models.StatusVal) ([]Order, error) {
	var orders []Order

	tx := s.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, `
		SELECT O.* FROM orders O JOIN statuses S ON O.id = S.order_id 
		           WHERE O.customer_id = $1 AND S."value" = $2 
		           ORDER BY S."time"
	`, customerId, statusVal)

	return orders, err
}
