package storage

import (
	"errors"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

var ErrOrderNotFound = errors.New("order is not found")

//	type Storage interface {
//		SetOrder(order models.Order) error
//		GetOrder(orderId models.IDType) (models.Order, error)
//		RemoveOrder(orderId models.IDType) error
//		GetOrderIDs() ([]models.IDType, error)
//		GetReturnIDs() ([]models.IDType, error)
//	}
//
//go:generate mkdir -p mock
//go:generate minimock -o ./mock/ -s .go -g
type Storage interface {
	CreateOrder(order models.Order) error

	GetOrder(orderId models.IDType) (models.Order, error)

	DeleteOrder(orderId models.IDType) error

	ChangeOrderStatus(orderId models.IDType, val models.Status) error

	GetCustomerOrdersWithStatus(customerId models.IDType, status models.StatusVal) ([]models.Order, error)

	GetOrderStatus(orderId models.IDType) (models.Status, error)

	GetReturnIDs() ([]models.IDType, error)
}
