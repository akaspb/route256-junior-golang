package storage

import (
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"os"
)

type JsonStorage struct {
	OrderStorage map[models.IDType]models.Order
	Returns      map[models.IDType]struct{}
	Path         string
}

type OrdersStorage struct {
	Orders  []models.Order  `json:"orders"`
	Returns []models.IDType `json:"returns"`
}

func InitJsonStorage(jsonPath string) (*JsonStorage, error) {
	jsonStorage := &JsonStorage{
		OrderStorage: nil,
		Path:         jsonPath,
		Returns:      nil,
	}

	if err := jsonStorage.readDataFromFile(); err != nil {
		return nil, err
	}

	return jsonStorage, nil
}

func (s *JsonStorage) SetOrder(order models.Order) error {
	if order.Status.Val == models.StatusReturn {
		if prevOrder, ok := s.OrderStorage[order.ID]; ok {
			if prevOrder.Status.Val != models.StatusReturn {
				s.Returns[order.ID] = struct{}{}
			}
		} else {
			s.Returns[order.ID] = struct{}{}
		}
	}

	s.OrderStorage[order.ID] = order

	if err := s.writeDataToFile(); err != nil {
		return err
	}

	return nil
}

func (s *JsonStorage) GetOrder(orderId models.IDType) (models.Order, error) {
	order, ok := s.OrderStorage[orderId]
	if !ok {
		return models.Order{}, ErrOrderNotFound
	}

	return order, nil
}

func (s *JsonStorage) GetOrderIDs() ([]models.IDType, error) {
	orderIDs := make([]models.IDType, 0, len(s.OrderStorage))
	for orderID, _ := range s.OrderStorage {
		orderIDs = append(orderIDs, orderID)
	}

	return orderIDs, nil
}

func (s *JsonStorage) RemoveOrder(orderID models.IDType) error {
	if _, ok := s.OrderStorage[orderID]; !ok {
		return fmt.Errorf("no order with orderId==%v in storage", orderID)
	}

	delete(s.OrderStorage, orderID)
	delete(s.Returns, orderID)

	if err := s.writeDataToFile(); err != nil {
		return err
	}

	return nil
}

func (s *JsonStorage) GetReturnIDs() ([]models.IDType, error) {
	returnIDs := make([]models.IDType, 0, len(s.Returns))
	for orderID, _ := range s.Returns {
		returnIDs = append(returnIDs, orderID)
	}

	return returnIDs, nil
}

func (s *JsonStorage) readDataFromFile() error {

	var file *os.File
	file, err := os.OpenFile(s.Path, os.O_RDWR, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.OpenFile(s.Path, os.O_CREATE, 0666)
			if err != nil {
				return err
			}

			s.OrderStorage = make(map[models.IDType]models.Order)
			s.Returns = make(map[models.IDType]struct{})
			return s.writeDataToFile()
		}
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

	s.Returns = make(map[models.IDType]struct{}, len(ordersStorage.Returns))
	for _, oderID := range ordersStorage.Returns {
		s.Returns[oderID] = struct{}{}
	}

	return nil
}

func (s *JsonStorage) writeDataToFile() error {
	var file *os.File
	var err error
	file, err = os.OpenFile(s.Path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	orders := make([]models.Order, 0, len(s.OrderStorage))
	for _, order := range s.OrderStorage {
		orders = append(orders, order)
	}

	returnIDs := make([]models.IDType, 0, len(s.Returns))
	for orderID, _ := range s.Returns {
		returnIDs = append(returnIDs, orderID)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(
		OrdersStorage{
			Orders:  orders,
			Returns: returnIDs,
		})
	if err != nil {
		return err
	}

	return nil
}
