package storage

import (
	"context"
	"errors"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage/postgres"
)

var (
	ErrOrderNotFound     = errors.New("order is not found")
	ErrOrderWithIdExists = errors.New("order with such id exist")
	ErrStatusNotFound    = errors.New("status is not found")
)

type Facade interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrder(ctx context.Context, orderID models.IDType) (models.Order, error)
	DeleteOrder(ctx context.Context, orderID models.IDType) error
	ChangeOrderStatus(ctx context.Context, orderID models.IDType, val models.Status) error
	GetCustomerOrdersWithStatus(ctx context.Context, customerID models.IDType, statusVal models.StatusVal) ([]models.Order, error)
	GetOrderStatus(ctx context.Context, orderID models.IDType) (models.Status, error)
	GetOrderIDsWhereStatus(ctx context.Context, statusVal models.StatusVal) ([]models.IDType, error)
}

type storageFacade struct {
	txManager    postgres.TransactionManager
	pgRepository *postgres.PgStorage
}

func NewStorageFacade(
	txManager postgres.TransactionManager,
	pgRepository *postgres.PgStorage,
) *storageFacade {
	return &storageFacade{
		txManager:    txManager,
		pgRepository: pgRepository,
	}
}

func (s *storageFacade) CreateOrder(ctx context.Context, order models.Order) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
		var err error
		err = s.pgRepository.CreateOrder(ctx, postgres.Order{
			ID:         order.ID,
			CustomerID: order.CustomerID,
			Expiry:     order.Expiry,
			Weight:     order.Weight,
			Cost:       order.Cost,
		})
		if err != nil {
			return err
		}

		err = s.pgRepository.CreateStatus(ctx, postgres.Status{
			Value:   models.StatusToStorage,
			Time:    order.Status.Time,
			OrderID: order.ID,
		})
		if err != nil {
			return err
		}

		if order.Pack != nil {
			pack := postgres.Pack{
				OrderID:        order.ID,
				Name:           order.Pack.Name,
				Cost:           order.Pack.Cost,
				MaxOrderWeight: order.Pack.MaxOrderWeight,
			}
			err = s.pgRepository.CreatePack(ctx, pack)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *storageFacade) GetOrder(ctx context.Context, orderID models.IDType) (models.Order, error) {
	storageOrder, err := s.pgRepository.GetOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, postgres.ErrorOrderNotFound) {
			return models.Order{}, ErrOrderNotFound
		}
		return models.Order{}, err
	}

	storageStatus, err := s.pgRepository.GetStatus(ctx, orderID)
	if err != nil {
		if errors.Is(err, postgres.ErrorStatusNotFound) {
			return models.Order{}, ErrOrderNotFound
		}
		return models.Order{}, err
	}

	var packPtr *models.Pack
	storagePack, err := s.pgRepository.GetPack(ctx, orderID)
	if err != nil {
		if !errors.Is(err, postgres.ErrorPackNotFound) {
			return models.Order{}, err
		}
	} else {
		packPtr = &models.Pack{
			Name:           storagePack.Name,
			Cost:           storagePack.Cost,
			MaxOrderWeight: storagePack.MaxOrderWeight,
		}
	}

	return models.Order{
		ID:         storageOrder.ID,
		CustomerID: storageOrder.CustomerID,
		Expiry:     storageOrder.Expiry,
		Weight:     storageOrder.Weight,
		Cost:       storageOrder.Cost,
		Pack:       packPtr,
		Status: models.Status{
			Value: storageStatus.Value,
			Time:  storageStatus.Time,
		},
	}, nil
}

func (s *storageFacade) DeleteOrder(ctx context.Context, orderId models.IDType) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
		if err := s.pgRepository.DeleteOrder(ctx, orderId); err != nil {
			return err
		}

		if err := s.pgRepository.DeletePack(ctx, orderId); !errors.Is(err, postgres.ErrorPackNotFound) {
			return err
		}

		return nil
	})
}

func (s *storageFacade) ChangeOrderStatus(ctx context.Context, orderId models.IDType, status models.Status) error {
	return s.pgRepository.SetStatus(ctx, postgres.Status{
		OrderID: orderId,
		Value:   status.Value,
		Time:    status.Time,
	})
}

func (s *storageFacade) GetCustomerOrdersWithStatus(ctx context.Context, customerId models.IDType, statusVal models.StatusVal) ([]models.Order, error) {
	var orders []models.Order

	storageOrders, err := s.pgRepository.GetCustomerOrdersWithStatus(ctx, customerId, statusVal)
	if err != nil {
		return nil, err
	}

	orders = make([]models.Order, len(storageOrders))
	for i, storageOrder := range storageOrders {
		var storagePack postgres.Pack
		var packPtr *models.Pack
		storagePack, err = s.pgRepository.GetPack(ctx, storageOrder.ID)
		if err != nil {
			if !errors.Is(err, postgres.ErrorPackNotFound) {
				return nil, err
			}
			packPtr = &models.Pack{
				Name:           storagePack.Name,
				Cost:           storagePack.Cost,
				MaxOrderWeight: storagePack.MaxOrderWeight,
			}
		}
		orders[i] = models.Order{
			ID:         storageOrder.ID,
			CustomerID: storageOrder.CustomerID,
			Expiry:     storageOrder.Expiry,
			Weight:     storageOrder.Weight,
			Cost:       storageOrder.Cost,
			Pack:       packPtr,
			Status:     models.Status{},
		}
	}

	return orders, err
}

func (s *storageFacade) GetOrderStatus(ctx context.Context, orderId models.IDType) (models.Status, error) {
	storageStatus, err := s.pgRepository.GetStatus(ctx, orderId)
	if err != nil {
		return models.Status{}, err
	}

	return models.Status{
		Value: storageStatus.Value,
		Time:  storageStatus.Time,
	}, nil
}

func (s *storageFacade) GetOrderIDsWhereStatus(ctx context.Context, statusVal models.StatusVal) ([]models.IDType, error) {
	var orderIDs []models.IDType

	err := s.txManager.RunReadUncommitted(ctx, func(ctxTx context.Context) error {
		var err error

		orderIDs, err = s.pgRepository.GetOrderIDsWhereStatus(ctx, statusVal)
		if err != nil {
			return err
		}

		return nil
	})

	return orderIDs, err
}
