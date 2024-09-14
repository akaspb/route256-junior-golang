package storage

import (
	"errors"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

var ErrOrderNotFound = errors.New("order is not found")

type Storage interface {
	SetOrder(order models.Order) error
	GetOrder(orderId models.IDType) (models.Order, error)
	RemoveOrder(orderId models.IDType) error
	GetOrderIDs() ([]models.IDType, error)
	GetReturnIDs() ([]models.IDType, error)
}
