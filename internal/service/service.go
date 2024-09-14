package service

import (
	"errors"
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"sort"
	"time"
)

const Day = 24 * time.Hour
const MaxReturnTime = 2 * Day

type OrderIDWithMsg struct {
	ID  models.IDType
	Msg string
	Ok  bool
}

type OrderIDWithExpiryAndStatus struct {
	ID      models.IDType
	Expiry  time.Time
	Expired bool
}

type ReturnOrderAndCustomer struct {
	OrderID    models.IDType
	CustomerID models.IDType
}

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
	_, err := s.orderStorage.GetOrder(orderID)

	if err == nil {
		return fmt.Errorf("order with ID==%v was accepted earlier", orderID)
	}

	if !errors.Is(err, storage.ErrOrderNotFound) {
		return err
	}

	if !isLessOrEqualTime(s.GetCurrentTime(), orderExpiry) {
		return errors.New("order expiry time can't be before current time")
	}

	return s.orderStorage.SetOrder(models.Order{
		ID:         orderID,
		CustomerID: customerID,
		Expiry:     orderExpiry,
		Status: models.Status{
			Val:  models.StatusToStorage,
			Time: s.GetCurrentTime(),
		},
	})
}

func (s *Service) ReturnOrder(orderID models.IDType) error {
	order, err := s.orderStorage.GetOrder(orderID)

	if err != nil {
		if errors.Is(err, storage.ErrOrderNotFound) {
			return fmt.Errorf("order with ID==%v is not in PVZ or has already been returned and given to courier", orderID)
		}
		return err
	}

	switch order.Status.Val {
	case models.StatusToCustomer:
		return fmt.Errorf("order with orderId==%v was taken by customer", orderID)
	case models.StatusReturn:
		break
	case models.StatusToStorage:
		if !isLessOrEqualTime(order.Expiry, s.GetCurrentTime()) {
			return fmt.Errorf("order with orderId==%v expiry time is not reached", orderID)
		}
	default: // if some unknown status
		return errors.New("unhandled error")
	}

	err = s.orderStorage.RemoveOrder(orderID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GiveOrderToCustomer(orderIDs []models.IDType, customerID models.IDType) ([]OrderIDWithMsg, error) {
	ordersToGive := make([]OrderIDWithMsg, 0)
	currTime := s.GetCurrentTime()

	for _, orderID := range orderIDs {
		order, err := s.orderStorage.GetOrder(orderID)
		if err != nil {
			if errors.Is(err, storage.ErrOrderNotFound) {
				// if order has not been delivered to storage yet
				ordersToGive = append(ordersToGive, OrderIDWithMsg{
					ID:  orderID,
					Msg: "Order has not been delivered to PVZ yet or was returned and given to courier",
					Ok:  false,
				})
				continue
			}
			return nil, err
		}

		// if order's customer is other person return error immediately
		if order.CustomerID != customerID {
			return nil, fmt.Errorf(
				"customer %v can't get the information about order with ID==%v",
				customerID, order.ID,
			)
		}

		switch order.Status.Val {
		case models.StatusToCustomer:
			ordersToGive = append(ordersToGive, OrderIDWithMsg{
				ID:  orderID,
				Msg: "Order was given to customer earlier",
				Ok:  false,
			})
			continue
		case models.StatusReturn:
			ordersToGive = append(ordersToGive, OrderIDWithMsg{
				ID:  orderID,
				Msg: "Order was returned to PVZ by customer",
				Ok:  false,
			})
			continue
		case models.StatusToStorage:
			if isLessOrEqualTime(currTime, order.Expiry) {
				ordersToGive = append(ordersToGive, OrderIDWithMsg{
					ID:  order.ID,
					Msg: "Give order to customer",
					Ok:  true,
				})
			} else {
				ordersToGive = append(ordersToGive, OrderIDWithMsg{
					ID:  order.ID,
					Msg: "Order expiry time was reached",
					Ok:  false,
				})
			}
		default: // if some unknown status
			return nil, errors.New("unhandled error")
		}
	}

	// add this loop with purpose if error occurred in previous loop no change in DB were made
	for _, order := range ordersToGive {
		if order.Ok {
			orderPtr, err := s.orderStorage.GetOrder(order.ID)
			if err != nil {
				return nil, err
			}

			orderPtr.Status.Val = models.StatusToCustomer
			orderPtr.Status.Time = currTime
			err = s.orderStorage.SetOrder(orderPtr)
			if err != nil {
				return nil, err
			}
		}
	}

	return ordersToGive, nil
}

func (s *Service) GetCustomerOrders(customerID models.IDType, n uint) ([]OrderIDWithExpiryAndStatus, error) {
	userOrders := make([]models.Order, 0)
	allOrdersIDs, err := s.orderStorage.GetOrderIDs()
	if err != nil {
		return nil, err
	}

	for _, orderID := range allOrdersIDs {
		order, err := s.orderStorage.GetOrder(orderID)
		if err != nil {
			if errors.Is(err, storage.ErrOrderNotFound) {
				return nil, errors.New("unhandled error")
			}
			return nil, err
		}

		if order.CustomerID == customerID {
			if order.Status.Val == models.StatusToStorage {
				userOrders = append(userOrders, order)
			}
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
	res := make([]OrderIDWithExpiryAndStatus, n)
	for i, userOrder := range userOrders {
		res[i] = OrderIDWithExpiryAndStatus{
			ID:      userOrder.ID,
			Expiry:  userOrder.Expiry,
			Expired: !isLessOrEqualTime(s.GetCurrentTime(), userOrder.Expiry),
		}
	}

	return res, nil
}

func (s *Service) ReturnOrderFromCustomer(customerID, orderID models.IDType) error {
	currTime := s.GetCurrentTime()

	order, err := s.orderStorage.GetOrder(orderID)

	if err != nil {
		if errors.Is(err, storage.ErrOrderNotFound) {
			return fmt.Errorf("order with ID==%v is not in PVZ or has already been returned", orderID)
		}
		return err
	}

	if order.CustomerID != customerID {
		return fmt.Errorf(
			"order with ID==%v can't be accepted for return from other customer %v",
			order.ID, customerID,
		)
	}

	if order.Status.Val == models.StatusReturn {
		return fmt.Errorf(
			"order with ID==%v was returned by customer %v",
			orderID, customerID,
		)
	}

	if order.Status.Val != models.StatusToCustomer {
		return fmt.Errorf(
			"order with ID==%v was not given to customer %v",
			orderID, customerID,
		)
	}

	if !isLessOrEqualTime(currTime, order.Status.Time.Add(MaxReturnTime)) {
		return fmt.Errorf("order with ID==%v return time elapsed", order.ID)
	}

	order.Status.Val = models.StatusReturn
	order.Status.Time = currTime

	err = s.orderStorage.SetOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetReturnsList(offset, limit int) ([]ReturnOrderAndCustomer, error) {
	orderIDsToReturn := make([]ReturnOrderAndCustomer, 0)
	returnIDs, err := s.orderStorage.GetReturnIDs()
	if err != nil {
		return nil, err
	}

	count := -offset
	for _, orderID := range returnIDs {
		order, err := s.orderStorage.GetOrder(orderID)
		if err != nil {
			return nil, err
		}

		if count >= 0 {
			orderIDsToReturn = append(orderIDsToReturn, ReturnOrderAndCustomer{
				OrderID:    order.ID,
				CustomerID: order.CustomerID,
			})
		}
		count++
		if count == limit {
			break
		}
	}

	if count < 1 {
		return nil, errors.New("offset value is too small")
	}

	return orderIDsToReturn, nil
}

func (s *Service) GetCurrentTime() time.Time {
	return s.startTime.Add(time.Now().Sub(s.systemStartTime))
}

func isLessOrEqualTime(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.Before(t2)
}
