package storage

import (
	"errors"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

var ErrOrderNotFound = errors.New("order is not found")

//go:generate mkdir -p mock
//go:generate minimock -o ./mock/ -s .go -g
type Storage interface {
	SetOrder(order models.Order) error
	GetOrder(orderId models.IDType) (models.Order, error)
	RemoveOrder(orderId models.IDType) error
	GetOrderIDs() ([]models.IDType, error)
	GetReturnIDs() ([]models.IDType, error)
}
