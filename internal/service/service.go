package service

import (
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"sort"
	"time"
)

const Day = time.Hour * 24
const ReturnTime = Day * 2

type Service struct {
	startTime       time.Time
	systemStartTime time.Time
	orderStorage    storage.Storage
}

func NewService(orderStorage storage.Storage, startTime time.Time, systemStartTime time.Time) *Service {
	return &Service{
		startTime:       startTime,
		systemStartTime: systemStartTime,
		orderStorage:    orderStorage,
	}
}

func (s *Service) AcceptOrderFromCourier(orderID, customerID models.IDType, orderExpiry time.Time) error {
	order, err := s.orderStorage.GetOrder(orderID)
	if err != nil {
		return err
	}

	if order != nil {
		return fmt.Errorf("order with ID==%v was accepted earlier", orderID)
	}

	if !isLessOrEqualTime(s.getCurrentTime(), orderExpiry) {
		return fmt.Errorf("order expiry time can't be before current time")
	}

	return s.orderStorage.SetOrder(models.Order{
		ID:     orderID,
		Expiry: orderExpiry,
		Status: models.Status{
			Val:        models.StatusToStorage,
			CustomerID: customerID,
			Time:       s.getCurrentTime(),
		},
	})
}

func (s *Service) ReturnOrder(orderID models.IDType) error {
	var order *models.Order
	order, err := s.orderStorage.GetOrder(orderID)
	if order == nil {
		return fmt.Errorf("order with ID==%v was't accepted earlier", orderID)
	}

	if order.Status.Val == models.StatusToStorage {
		return fmt.Errorf("order with orderId==%v is in storage yet", orderID)
	}

	if order.Status.Val == models.StatusReturn {
		return fmt.Errorf("order with orderId==%v was already returned by customer", orderID)
	}

	if !isLessOrEqualTime(order.Expiry, s.getCurrentTime()) {
		return fmt.Errorf("order with orderId==%v expiry time is not reached", orderID)
	}

	order.Status.Val = models.StatusReturn
	order.Status.Time = s.getCurrentTime()
	err = s.orderStorage.SetOrder(*order)
	if err != nil {
		return err
	}

	//err = s.orderStorage.RemoveOrder(orderID)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (s *Service) GiveOrderToCustomer(orderIDs []models.IDType, customerID models.IDType) ([]models.IDType, error) {
	orderIDsToGive := make([]models.IDType, 0)
	currTime := s.getCurrentTime()

	for _, orderID := range orderIDs {
		var order *models.Order
		order, err := s.orderStorage.GetOrder(orderID)
		if err != nil {
			return nil, err
		}

		if order == nil {
			continue // order was not delivered to storage yet
		}

		if order.CustomerID != customerID {
			return nil, fmt.Errorf(
				"order (ID=%v, CustomerID=%v) can't be given to customer %v",
				order.ID, order.CustomerID, customerID,
			)
		}

		if order.Status.Val == models.StatusToStorage {
			if isLess(currTime, order.Expiry) {
				orderIDsToGive = append(orderIDsToGive, order.ID)
			}
		} else {
			// no order with this id in storage right now
		}
	}

	for _, orderID := range orderIDsToGive {
		order, err := s.orderStorage.GetOrder(orderID)
		if err != nil {
			return nil, err
		}

		order.Status.Val = models.StatusToCustomer
		order.Status.Time = currTime
		err = s.orderStorage.SetOrder(*order)
		if err != nil {
			return nil, err
		}
	}

	return orderIDsToGive, nil
}

func (s *Service) GetCustomerOrders(customerID models.IDType, n uint) ([]models.IDType, error) {
	userOrders := make([]models.Order, 0)
	allOrdersIDs, err := s.orderStorage.GetOrderIDs()
	if err != nil {
		return nil, err
	}

	for _, orderID := range allOrdersIDs {
		order, err := s.orderStorage.GetOrder(orderID)
		if err != nil {
			return nil, err
		}

		if order == nil {
			return nil, fmt.Errorf("unhandled error")
		}

		if order.Status.Val == models.StatusToStorage && order.CustomerID == customerID {
			userOrders = append(userOrders, *order)
		}
	}

	sort.Slice(userOrders, func(i, j int) bool {
		if userOrders[i].Status.Time.Equal(userOrders[j].Status.Time) {
			return userOrders[i].ID < userOrders[j].ID
		}
		return userOrders[i].Status.Time.Before(userOrders[j].Status.Time)
	})

	if n == 0 {
		n = uint(len(userOrders))
	}

	n = min(n, uint(len(userOrders)))

	userOrders = userOrders[len(userOrders)-int(n):]
	userOrderIDs := make([]models.IDType, n)
	for i, userOrder := range userOrders {
		userOrderIDs[i] = userOrder.ID
	}

	return userOrderIDs, nil
}

func (s *Service) ReturnOrderFromCustomer(customerID, orderID models.IDType) error {
	currTime := s.getCurrentTime()

	order, err := s.orderStorage.GetOrder(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return fmt.Errorf("order with ID==%v was't accepted earlier", orderID)
	}

	if order.Status.Val != models.StatusToCustomer {
		return fmt.Errorf(
			"order with ID==%v was not given to customer %v",
			orderID, customerID,
		)
	}

	if order.CustomerID != customerID {
		return fmt.Errorf(
			"order (ID=%v, CustomerID=%v) can't be accepted for return from customer %v",
			order.ID, order.CustomerID, customerID,
		)
	}

	if !isLess(currTime, order.Status.Time.Add(ReturnTime)) {
		return fmt.Errorf("order with ID==%v return time elapsed")
	}

	order.Status.Val = models.StatusReturn
	order.Status.Time = currTime

	err = s.orderStorage.SetOrder(*order)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetReturnsList() ([]models.IDType, error) {
	orderIDsToReturn := make([]models.IDType, 0)
	allOrdersIDs, err := s.orderStorage.GetOrderIDs()
	if err != nil {
		return nil, err
	}

	for _, orderID := range allOrdersIDs {
		order, err := s.orderStorage.GetOrder(orderID)
		if err != nil {
			return nil, err
		}

		if order == nil {
			return nil, fmt.Errorf("unhandled error")
		}

		if order.Status.Val == models.StatusReturn {
			orderIDsToReturn = append(orderIDsToReturn, order.ID)
		}
	}

	return orderIDsToReturn, nil
}

func (s *Service) getCurrentTime() time.Time {
	return s.startTime.Add(time.Now().Sub(s.systemStartTime))
}

func isLess(t1, t2 time.Time) bool {
	return t1.Before(t2)
}

func isLessOrEqualTime(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.Before(t2)
}
