package storage

import (
	"context"
	"errors"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"sync"
)

var (
	ErrorContextDone = errors.New("context was done")
)

type Storage struct {
	orders sync.Map
}

// NewStorage returns test in memory Storage (for testing only, maybe not optimized)
func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) CreateOrder(ctx context.Context, order models.Order) error {
	select {
	case <-ctx.Done():
		return ErrorContextDone
	default:
	}

	if _, ok := s.orders.Load(order.ID); ok {
		return storage.ErrOrderWithIdExists
	}

	s.orders.Store(order.ID, order)

	return nil
}

func (s *Storage) GetOrder(ctx context.Context, orderID models.IDType) (models.Order, error) {
	select {
	case <-ctx.Done():
		return models.Order{}, ErrorContextDone
	default:
	}

	if order, ok := s.orders.Load(orderID); ok {
		return order.(models.Order), nil
	}

	return models.Order{}, storage.ErrOrderNotFound
}

func (s *Storage) DeleteOrder(ctx context.Context, orderID models.IDType) error {
	select {
	case <-ctx.Done():
		return ErrorContextDone
	default:
	}

	if _, ok := s.orders.Load(orderID); !ok {
		return storage.ErrOrderNotFound
	}

	s.orders.Delete(orderID)

	return nil
}

func (s *Storage) ChangeOrderStatus(ctx context.Context, orderID models.IDType, status models.Status) error {
	select {
	case <-ctx.Done():
		return ErrorContextDone
	default:
	}

	if order, ok := s.orders.Load(orderID); ok {
		order := order.(models.Order)
		order.Status = status
		s.orders.Store(orderID, order)
		return nil
	}

	return storage.ErrOrderNotFound
}

func (s *Storage) GetCustomerOrdersWithStatus(ctx context.Context, customerID models.IDType, statusVal models.StatusVal) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ErrorContextDone
	default:
	}

	orders := make([]models.Order, 0)
	s.orders.Range(func(_, orderAny interface{}) bool {
		order := orderAny.(models.Order)
		if order.CustomerID == customerID && order.Status.Value == statusVal {
			orders = append(orders, order)
		}
		return true
	})

	return orders, nil
}

func (s *Storage) GetOrderStatus(ctx context.Context, orderID models.IDType) (models.Status, error) {
	select {
	case <-ctx.Done():
		return models.Status{}, ErrorContextDone
	default:
	}

	if status, ok := s.orders.Load(orderID); ok {
		return status.(models.Status), nil
	}

	return models.Status{}, storage.ErrOrderNotFound
}

func (s *Storage) GetOrderIDsWhereStatus(ctx context.Context, statusVal models.StatusVal, offset, limit uint) ([]models.IDType, error) {
	select {
	case <-ctx.Done():
		return nil, ErrorContextDone
	default:
	}

	orderIDs := make([]models.IDType, 0)
	s.orders.Range(func(_, orderAny interface{}) bool {
		order := orderAny.(models.Order)
		if order.Status.Value == statusVal {
			if offset > 0 {
				offset--
			} else {
				if limit > 0 {
					orderIDs = append(orderIDs, order.ID)
					limit--
				} else {
					return false
				}
			}
		}
		return true
	})

	return orderIDs, nil
}

func (s *Storage) GetOrderCustomerID(ctx context.Context, orderID models.IDType) (models.IDType, error) {
	select {
	case <-ctx.Done():
		return 0, ErrorContextDone
	default:
	}

	if order, ok := s.orders.Load(orderID); ok {
		return order.(models.Order).CustomerID, nil
	}

	return 0, storage.ErrOrderNotFound
}

func (s *Storage) FillWithOrders(ctx context.Context, orders ...models.Order) error {
	select {
	case <-ctx.Done():
		return ErrorContextDone
	default:
	}

	for _, order := range orders {
		if err := s.CreateOrder(ctx, order); err != nil {
			return err
		}
	}

	return nil
}
