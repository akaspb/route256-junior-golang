package storage

import (
	"context"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
)

type Storage struct {
	orders map[models.IDType]models.Order
}

// NewStorage returns test in memory Storage (for testing only, maybe not optimized)
func NewStorage() *Storage {
	return &Storage{orders: make(map[models.IDType]models.Order)}
}

func (s *Storage) CreateOrder(ctx context.Context, order models.Order) error {
	if _, ok := s.orders[order.ID]; ok {
		return storage.ErrOrderWithIdExists
	}

	s.orders[order.ID] = order

	return nil
}

func (s *Storage) GetOrder(ctx context.Context, orderID models.IDType) (models.Order, error) {
	if order, ok := s.orders[orderID]; ok {
		return order, nil
	}

	return models.Order{}, storage.ErrOrderNotFound
}

func (s *Storage) DeleteOrder(ctx context.Context, orderID models.IDType) error {
	if _, ok := s.orders[orderID]; !ok {
		return storage.ErrOrderNotFound
	}

	delete(s.orders, orderID)

	return nil
}

func (s *Storage) ChangeOrderStatus(ctx context.Context, orderID models.IDType, status models.Status) error {
	if _, ok := s.orders[orderID]; !ok {
		return storage.ErrOrderNotFound
	}

	order := s.orders[orderID]
	order.Status = status
	s.orders[orderID] = order

	return nil
}

func (s *Storage) GetCustomerOrdersWithStatus(ctx context.Context, customerID models.IDType, statusVal models.StatusVal) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	for _, order := range s.orders {
		if order.CustomerID == customerID && order.Status.Value == statusVal {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (s *Storage) GetOrderStatus(ctx context.Context, orderID models.IDType) (models.Status, error) {
	if _, ok := s.orders[orderID]; !ok {
		return models.Status{}, storage.ErrOrderNotFound
	}

	return s.orders[orderID].Status, nil
}

func (s *Storage) GetOrderIDsWhereStatus(ctx context.Context, statusVal models.StatusVal, offset, limit uint) ([]models.IDType, error) {
	orderIDs := make([]models.IDType, 0)
	for _, order := range s.orders {
		if order.Status.Value == statusVal {
			if offset > 0 {
				offset--
			} else {
				if limit > 0 {
					orderIDs = append(orderIDs, order.ID)
					limit--
				} else {
					break
				}
			}
		}
	}

	return orderIDs, nil
}

func (s *Storage) GetOrderCustomerID(ctx context.Context, orderID models.IDType) (models.IDType, error) {
	if order, ok := s.orders[orderID]; ok {
		return order.CustomerID, nil
	}

	return 0, storage.ErrOrderNotFound
}

func (s *Storage) FillWithOrders(ctx context.Context, orders ...models.Order) error {
	for _, order := range orders {
		if err := s.CreateOrder(ctx, order); err != nil {
			return err
		}
	}

	return nil
}
