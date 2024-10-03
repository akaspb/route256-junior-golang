package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
)

const Day = 24 * time.Hour
const MaxReturnTime = 2 * Day

var (
	ErrorCustomerID              = errors.New("operation is forbidden for not this order customer")
	ErrorOrderWasAccepted        = errors.New("order was accepted earlier")
	ErrorOrderExpiredAlready     = errors.New("order expired")
	ErrorOderNegativeCost        = errors.New("order cost can't be negative")
	ErrorOderNegativeWeight      = errors.New("order weight can't be negative")
	ErrorOrderWasTakenByCustomer = errors.New("order was taken by customer")
	ErrorOrderWasNotFounded      = errors.New("order is not in PVZ or has already been returned and given to courier")
)

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
	orderStorage    storage.Facade
}

func NewService(
	orderStorage storage.Facade,
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
	ctx context.Context,
	acceptOrderDTO AcceptOrderDTO,
) error {
	orderID := acceptOrderDTO.OrderID
	orderCost := acceptOrderDTO.OrderCost
	oderWeight := acceptOrderDTO.OderWeight
	customerID := acceptOrderDTO.CustomerID
	pack := acceptOrderDTO.Pack
	orderExpiry := acceptOrderDTO.OrderExpiry

	_, err := s.orderStorage.GetOrder(ctx, orderID)

	if err == nil {
		return ErrorOrderWasAccepted
	}

	if !errors.Is(err, storage.ErrOrderNotFound) {
		return err
	}

	if !isLessOrEqualTime(s.GetCurrentTime(), orderExpiry) {
		return ErrorOrderExpiredAlready
	}

	if orderCost < 0 {
		return ErrorOderNegativeCost
	}

	if oderWeight < 0 {
		return ErrorOderNegativeWeight
	}

	if pack != nil {
		packCost, err := s.Packaging.PackOrder(*pack, oderWeight)
		if err != nil {
			return err
		}
		orderCost += packCost
	}

	return s.orderStorage.CreateOrder(ctx, models.Order{
		ID:         orderID,
		CustomerID: customerID,
		Weight:     oderWeight,
		Cost:       orderCost,
		Expiry:     orderExpiry,
		Pack:       pack,
		Status:     models.Status{ChangedAt: s.GetCurrentTime()},
	})
}

func (s *Service) ReturnOrder(ctx context.Context, orderID models.IDType) error {
	order, err := s.orderStorage.GetOrder(ctx, orderID)

	if err != nil {
		if errors.Is(err, storage.ErrOrderNotFound) {
			return ErrorOrderWasNotFounded
		}
		return err
	}

	switch order.Status.Value {
	case models.StatusToCustomer:
		return ErrorOrderWasTakenByCustomer
	case models.StatusReturn:
		break
	case models.StatusToStorage:
		if !isLessOrEqualTime(order.Expiry, s.GetCurrentTime()) {
			return fmt.Errorf("order with orderId==%v expiry time is not reached", orderID)
		}
	default: // if some unknown status
		return errors.New("unhandled error")
	}

	err = s.orderStorage.DeleteOrder(ctx, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GiveOrderToCustomer(ctx context.Context, orderIDs []models.IDType, customerID models.IDType) ([]OrderIDWithMsg, error) {
	ordersWantedToBeGiven := make([]OrderIDWithMsg, 0)
	currTime := s.GetCurrentTime()

	for _, orderID := range orderIDs {
		order, err := s.orderStorage.GetOrder(ctx, orderID)

		if err != nil {
			if errors.Is(err, storage.ErrOrderNotFound) {
				// if order has not been delivered to storage yet
				ordersWantedToBeGiven = append(ordersWantedToBeGiven, OrderIDWithMsg{
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
			return nil, ErrorCustomerID
		}

		packagingName := ""
		if order.Pack != nil {
			packagingName = order.Pack.Name
		}

		switch order.Status.Value {
		case models.StatusToCustomer:
			ordersWantedToBeGiven = append(ordersWantedToBeGiven, OrderIDWithMsg{
				ID:      orderID,
				Cost:    order.Cost,
				Package: packagingName,
				Msg:     "Order was given to customer earlier",
				Ok:      false,
			})
			continue
		case models.StatusReturn:
			ordersWantedToBeGiven = append(ordersWantedToBeGiven, OrderIDWithMsg{
				ID:      orderID,
				Cost:    order.Cost,
				Package: packagingName,
				Msg:     "Order was returned to PVZ by customer",
				Ok:      false,
			})
			continue
		case models.StatusToStorage:
			if isLessOrEqualTime(currTime, order.Expiry) {
				ordersWantedToBeGiven = append(ordersWantedToBeGiven, OrderIDWithMsg{
					ID:      order.ID,
					Cost:    order.Cost,
					Package: packagingName,
					Msg:     "Give order to customer",
					Ok:      true,
				})
			} else {
				ordersWantedToBeGiven = append(ordersWantedToBeGiven, OrderIDWithMsg{
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
	for _, order := range ordersWantedToBeGiven {
		if order.Ok {
			if err := s.orderStorage.ChangeOrderStatus(ctx, order.ID, models.Status{
				Value:     models.StatusToCustomer,
				ChangedAt: currTime,
			}); err != nil {
				return nil, err
			}
		}
	}

	return ordersWantedToBeGiven, nil
}

func (s *Service) GetCustomerOrders(ctx context.Context, customerID models.IDType, n uint) ([]OrderIDWithExpiryAndStatus, error) {
	userOrders, err := s.orderStorage.GetCustomerOrdersWithStatus(ctx, customerID, models.StatusToStorage)
	if err != nil {
		return nil, err
	}

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

func (s *Service) ReturnOrderFromCustomer(ctx context.Context, customerID, orderID models.IDType) error {
	currTime := s.GetCurrentTime()

	order, err := s.orderStorage.GetOrder(ctx, orderID)

	if err != nil {
		if errors.Is(err, storage.ErrOrderNotFound) {
			return ErrorOrderWasNotFounded
		}
		return err
	}

	if order.CustomerID != customerID {
		return ErrorCustomerID
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

	if !isLessOrEqualTime(currTime, order.Status.ChangedAt.Add(MaxReturnTime)) {
		return ErrorOrderExpiredAlready
	}

	err = s.orderStorage.ChangeOrderStatus(ctx, order.ID, models.Status{
		Value:     models.StatusReturn,
		ChangedAt: currTime,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetReturnsList(ctx context.Context, offset, limit int) ([]ReturnOrderAndCustomer, error) {
	if offset < 0 {
		return nil, errors.New("offset value must be >= 0")
	}

	if limit <= 0 {
		return nil, errors.New("limit value must be > 0")
	}

	returnOrderIDs, err := s.orderStorage.GetOrderIDsWhereStatus(ctx, models.StatusReturn, uint(offset), uint(limit))
	if err != nil {
		return nil, err
	}

	orderIDsToReturn := make([]ReturnOrderAndCustomer, len(returnOrderIDs))
	for i, orderID := range returnOrderIDs {
		customerID, err := s.orderStorage.GetOrderCustomerID(ctx, orderID)
		if err != nil {
			return nil, err
		}

		orderIDsToReturn[i] = ReturnOrderAndCustomer{
			OrderID:    orderID,
			CustomerID: customerID,
		}
	}

	if len(orderIDsToReturn) < 1 {
		return nil, errors.New("no orders to show with such offset and limit")
	}

	return orderIDsToReturn, nil
}

func (s *Service) GetCurrentTime() time.Time {
	return s.startTime.Add(time.Now().Sub(s.systemStartTime))
}

func isLessOrEqualTime(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.Before(t2)
}

func (s *Service) SetStartTime(startTime time.Time) {
	s.startTime = startTime
}
