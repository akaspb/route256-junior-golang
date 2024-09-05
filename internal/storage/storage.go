package storage

import "gitlab.ozon.dev/siralexpeter/Homework/internal/models"

type Storage interface {
	SetOrder(order models.Order) error
	GetOrder(orderId models.IDType) (*models.Order, error)
	RemoveOrder(orderId models.IDType) error
	GetOrderIDs() ([]models.IDType, error)
}
