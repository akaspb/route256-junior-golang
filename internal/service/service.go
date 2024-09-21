package service

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
)

const Day = 24 * time.Hour
const MaxReturnTime = 2 * Day

type OrderIDWithMsg struct {
	ID      models.IDType
	Cost    models.CostType
	Package string
	Msg     string
	Ok      bool
}

type OrderIDWithExpiryAndStatus struct {
	ID      models.IDType
	Package string
	Cost    models.CostType
	Expiry  time.Time
	Expired bool
}

type ReturnOrderAndCustomer struct {
	OrderID    models.IDType
	CustomerID models.IDType
}

type Service struct {
	Packaging       *packaging.Packaging
	startTime       time.Time
	systemStartTime time.Time
	orderStorage    storage.Storage
}

func NewService(
	orderStorage storage.Storage,
	packagingSrvc *packaging.Packaging,
	startTime time.Time,
	systemStartTime time.Time,
) *Service {
	return &Service{
		Packaging:       packagingSrvc,
		startTime:       startTime,
		systemStartTime: systemStartTime,
		orderStorage:    orderStorage,
	}
}

type AcceptOrderDTO struct {
	OrderID     models.IDType
	OrderCost   models.CostType
	OderWeight  models.WeightType
	CustomerID  models.IDType
	Pack        *models.Pack
	OrderExpiry time.Time
}

func (s *Service) AcceptOrderFromCourier(
	acceptOrderDTO AcceptOrderDTO,
) error {
	orderID := acceptOrderDTO.OrderID
	orderCost := acceptOrderDTO.OrderCost
	oderWeight := acceptOrderDTO.OderWeight
	customerID := acceptOrderDTO.CustomerID
	pack := acceptOrderDTO.Pack
	orderExpiry := acceptOrderDTO.OrderExpiry

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

	if orderCost < 0 {
		return errors.New("order cost can't be negative")
	}

	if oderWeight < 0 {
		return errors.New("order weight can't be negative")
	}

	if pack != nil {
		packCost, err := s.Packaging.PackOrder(*pack, oderWeight)
		if err != nil {
			return err
		}
		orderCost += packCost
	}

	return s.orderStorage.SetOrder(models.Order{
		ID:         orderID,
		CustomerID: customerID,
		Weight:     oderWeight,
		Cost:       orderCost,
		Expiry:     orderExpiry,
		Pack:       pack,
		Status: models.Status{
			Value: models.StatusToStorage,
			Time:  s.GetCurrentTime(),
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

	switch order.Status.Value {
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
					ID:      orderID,
					Package: "",
					Msg:     "Order has not been delivered to PVZ yet or was returned and given to courier",
					Ok:      false,
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

		packagingName := ""
		if order.Pack != nil {
			packagingName = order.Pack.Name
		}

		switch order.Status.Value {
		case models.StatusToCustomer:
			ordersToGive = append(ordersToGive, OrderIDWithMsg{
				ID:      orderID,
				Cost:    order.Cost,
				Package: packagingName,
				Msg:     "Order was given to customer earlier",
				Ok:      false,
			})
			continue
		case models.StatusReturn:
			ordersToGive = append(ordersToGive, OrderIDWithMsg{
				ID:      orderID,
				Cost:    order.Cost,
				Package: packagingName,
				Msg:     "Order was returned to PVZ by customer",
				Ok:      false,
			})
			continue
		case models.StatusToStorage:
			if isLessOrEqualTime(currTime, order.Expiry) {
				ordersToGive = append(ordersToGive, OrderIDWithMsg{
					ID:      order.ID,
					Cost:    order.Cost,
					Package: packagingName,
					Msg:     "Give order to customer",
					Ok:      true,
				})
			} else {
				ordersToGive = append(ordersToGive, OrderIDWithMsg{
					ID:      order.ID,
					Cost:    order.Cost,
					Package: packagingName,
					Msg:     "Order expiry time was reached",
					Ok:      false,
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

			orderPtr.Status.Value = models.StatusToCustomer
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
			if order.Status.Value == models.StatusToStorage {
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
		packagingName := ""
		if userOrder.Pack != nil {
			packagingName = userOrder.Pack.Name
		}

		res[i] = OrderIDWithExpiryAndStatus{
			ID:      userOrder.ID,
			Cost:    userOrder.Cost,
			Package: packagingName,
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

	if order.Status.Value == models.StatusReturn {
		return fmt.Errorf(
			"order with ID==%v was returned by customer %v",
			orderID, customerID,
		)
	}

	if order.Status.Value != models.StatusToCustomer {
		return fmt.Errorf(
			"order with ID==%v was not given to customer %v",
			orderID, customerID,
		)
	}

	if !isLessOrEqualTime(currTime, order.Status.Time.Add(MaxReturnTime)) {
		return fmt.Errorf("order with ID==%v return time elapsed", order.ID)
	}

	order.Status.Value = models.StatusReturn
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
