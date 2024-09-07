package storage

import (
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"os"
)

type JsonStorage struct {
	OrderStorage map[models.IDType]models.Order
	Path         string
}

type OrdersStorage struct {
	Orders []models.Order `json:"orders"`
}

func InitJsonStorage(jsonPath string) (*JsonStorage, error) {
	jsonStorage := &JsonStorage{
		OrderStorage: nil,
		Path:         jsonPath,
	}

	if err := jsonStorage.readDataFromFile(); err != nil {
		return nil, err
	}

	return jsonStorage, nil
}

func (s *JsonStorage) SetOrder(order models.Order) error {
	if s.OrderStorage == nil {
		return fmt.Errorf("storage was not loaded")
	}

	s.OrderStorage[order.ID] = order

	if err := s.writeDataFromFile(); err != nil {
		return err
	}

	return nil
}

func (s *JsonStorage) GetOrder(orderId models.IDType) (*models.Order, error) {
	if s.OrderStorage == nil {
		return nil, fmt.Errorf("storage was not loaded")
	}

	order, ok := s.OrderStorage[orderId]
	if !ok {
		return nil, nil
	}

	return &order, nil
}

func (s *JsonStorage) GetOrderIDs() ([]models.IDType, error) {
	if s.OrderStorage == nil {
		return nil, fmt.Errorf("storage was not loaded")
	}

	orderIDs := make([]models.IDType, 0, len(s.OrderStorage))
	for orderID, _ := range s.OrderStorage {
		orderIDs = append(orderIDs, orderID)
	}

	return orderIDs, nil
}

func (s *JsonStorage) RemoveOrder(orderID models.IDType) error {
	if s.OrderStorage == nil {
		return fmt.Errorf("storage was not loaded")
	}

	if _, ok := s.OrderStorage[orderID]; !ok {
		return fmt.Errorf("no order with orderId==%v in storage", orderID)
	}

	delete(s.OrderStorage, orderID)

	if err := s.writeDataFromFile(); err != nil {
		return err
	}

	return nil
}

func (s *JsonStorage) readDataFromFile() error {
	var file *os.File
	file, err := os.OpenFile(s.Path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var ordersStorage OrdersStorage
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ordersStorage)
	if err != nil {
		return err
	}

	orders := ordersStorage.Orders

	s.OrderStorage = make(map[models.IDType]models.Order, len(orders))
	for _, order := range orders {
		s.OrderStorage[order.ID] = order
	}

	return nil
}

func (s *JsonStorage) writeDataFromFile() (err error) {
	var file *os.File
	file, err = os.OpenFile(s.Path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var orders []models.Order
	if len(s.OrderStorage) == 0 {
		orders = make([]models.Order, 0)
	} else {
		orders = make([]models.Order, 0, len(s.OrderStorage))
		for _, order := range s.OrderStorage {
			orders = append(orders, order)
		}
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(OrdersStorage{Orders: orders})
	if err != nil {
		return err
	}

	return nil
}
