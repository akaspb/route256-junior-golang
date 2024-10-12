package service

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"golang.org/x/sync/errgroup"
)

const (
	Day           = 24 * time.Hour
	MaxReturnTime = 2 * Day
)

type ArgumentError struct {
	text string
}

func NewArgumentError(text string) *ArgumentError {
	return &ArgumentError{text: text}
}

func (e *ArgumentError) Error() string {
	return e.text
}

var (
	ErrorCustomerID              = NewArgumentError("operation is forbidden for not this order customer")
	ErrorOrderWasAccepted        = NewArgumentError("order was accepted earlier")
	ErrorOrderExpiredAlready     = NewArgumentError("order expired")
	ErrorOderNegativeCost        = NewArgumentError("order cost can't be negative")
	ErrorOderNegativeWeight      = NewArgumentError("order weight can't be negative")
	ErrorOrderWasTakenByCustomer = NewArgumentError("order was taken by customer")
	ErrorOrderWasNotFounded      = NewArgumentError("order is not in PVZ or has already been returned and given to courier")
	ErrorInvalidWorkerCount      = NewArgumentError("worker count must be > 0 and <= max thread count")
	ErrorInvalidOffsetValue      = NewArgumentError("offset value must be >= 0")
	ErrorInvalidLimitValue       = NewArgumentError("limit value must be > 0")
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
	Weight  models.WeightType
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
	workerCount     int
}

func NewService(
	orderStorage storage.Facade,
	packagingSrvc *packaging.Packaging,
	startTime time.Time,
	systemStartTime time.Time,
	workerCount int,
) (*Service, error) {
	s := &Service{
		Packaging:       packagingSrvc,
		startTime:       startTime,
		systemStartTime: systemStartTime,
		orderStorage:    orderStorage,
	}
	if err := s.ChangeWorkerCount(workerCount); err != nil {
		return nil, err
	}

	return s, nil
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

func (s *Service) getCustomerOrderIDWithMsg(ctx context.Context, orderID, customerID models.IDType, currTime time.Time) (OrderIDWithMsg, error) {
	order, err := s.orderStorage.GetOrder(ctx, orderID)

	if err != nil {
		if errors.Is(err, storage.ErrOrderNotFound) {
			// if order has not been delivered to storage yet
			return OrderIDWithMsg{
				ID:  orderID,
				Msg: "Order has not been delivered to PVZ yet or was returned and given to courier",
				Ok:  false,
			}, nil
		}
		return OrderIDWithMsg{}, err
	}

	// if order's customer is other person return error immediately
	if order.CustomerID != customerID {
		return OrderIDWithMsg{}, ErrorCustomerID
	}

	packagingName := ""
	if order.Pack != nil {
		packagingName = order.Pack.Name
	}

	switch order.Status.Value {
	case models.StatusToCustomer:
		return OrderIDWithMsg{
			ID:      orderID,
			Cost:    order.Cost,
			Package: packagingName,
			Msg:     "Order was given to customer earlier",
			Ok:      false,
		}, nil
	case models.StatusReturn:
		return OrderIDWithMsg{
			ID:      orderID,
			Cost:    order.Cost,
			Package: packagingName,
			Msg:     "Order was returned to PVZ by customer",
			Ok:      false,
		}, nil
	case models.StatusToStorage:
		if isLessOrEqualTime(currTime, order.Expiry) {
			return OrderIDWithMsg{
				ID:      order.ID,
				Cost:    order.Cost,
				Package: packagingName,
				Msg:     "Give order to customer",
				Ok:      true,
			}, nil
		} else {
			return OrderIDWithMsg{
				ID:      order.ID,
				Cost:    order.Cost,
				Package: packagingName,
				Msg:     "Order expiry time was reached",
				Ok:      false,
			}, nil
		}
	default: // if some unknown status
		return OrderIDWithMsg{}, errors.New("unhandled error (unknown order status value)")
	}
}

func (s *Service) GiveOrderToCustomer(ctx context.Context, orderIDs []models.IDType, customerID models.IDType) ([]OrderIDWithMsg, error) {
	currTime := s.GetCurrentTime()

	// Worker Pool
	group1, ctx1 := errgroup.WithContext(ctx)
	orderIDsChan := make(chan models.IDType, len(orderIDs))
	orderIDWithMsgChan := make(chan OrderIDWithMsg, len(orderIDs))
	okOrderIDsChan := make(chan models.IDType, len(orderIDs))

	for i := 0; i < s.workerCount; i++ {
		group1.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case orderID, ok := <-orderIDsChan:
					if !ok {
						return nil
					}

					orderIDWithMsg, err := s.getCustomerOrderIDWithMsg(ctx1, orderID, customerID, currTime)
					if err != nil {
						return err
					}

					orderIDWithMsgChan <- orderIDWithMsg
					if orderIDWithMsg.Ok {
						okOrderIDsChan <- orderIDWithMsg.ID
					}
				}
			}
		})
	}

	for _, orderID := range orderIDs {
		orderIDsChan <- orderID
	}
	close(orderIDsChan)

	if err := group1.Wait(); err != nil {
		defer close(orderIDWithMsgChan)
		defer close(okOrderIDsChan)
		return nil, err
	}
	close(orderIDWithMsgChan)
	close(okOrderIDsChan)

	// add this goroutines with purpose if error occurred in previous goroutines no change in DB were made
	group2, ctx2 := errgroup.WithContext(ctx)

	for i := 0; i < s.workerCount; i++ {
		group2.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case orderID, ok := <-okOrderIDsChan:
					if !ok {
						return nil
					}

					if err := s.orderStorage.ChangeOrderStatus(ctx2, orderID, models.Status{
						Value:     models.StatusToCustomer,
						ChangedAt: currTime,
					}); err != nil {
						return err
					}
				}
			}
		})
	}

	if err := group2.Wait(); err != nil {
		return nil, err
	}

	res := make([]OrderIDWithMsg, 0, len(orderIDWithMsgChan))
	for orderIDWithMsg := range orderIDWithMsgChan {
		res = append(res, orderIDWithMsg)
	}

	return res, nil
}

func (s *Service) getOrderIDWithExpiryAndStatus(ctx context.Context, orderID models.IDType, currTime time.Time) (OrderIDWithExpiryAndStatus, error) {
	order, err := s.orderStorage.GetOrder(ctx, orderID)
	if err != nil {
		return OrderIDWithExpiryAndStatus{}, err
	}

	packagingName := ""
	if order.Pack != nil {
		packagingName = order.Pack.Name
	}

	return OrderIDWithExpiryAndStatus{
		ID:      order.ID,
		Cost:    order.Cost,
		Weight:  order.Weight,
		Package: packagingName,
		Expiry:  order.Expiry,
		Expired: !isLessOrEqualTime(currTime, order.Expiry),
	}, nil
}

func (s *Service) GetCustomerOrders(ctx context.Context, customerID models.IDType, n uint) ([]OrderIDWithExpiryAndStatus, error) {
	var orderIDs []models.IDType
	var err error
	currTime := s.GetCurrentTime()

	if n == 0 {
		orderIDs, err = s.orderStorage.GetCustomerOrderIDsWithStatus(ctx, customerID, models.StatusToStorage)
		if err != nil {
			return nil, err
		}
	} else {
		orderIDs, err = s.orderStorage.GetNCustomerOrderIDsWithStatus(ctx, customerID, models.StatusToStorage, n)
		if err != nil {
			return nil, err
		}
	}

	// Worker Pool
	group, ctx := errgroup.WithContext(ctx)
	orderIDsChan := make(chan models.IDType, len(orderIDs))
	resChan := make(chan OrderIDWithExpiryAndStatus, len(orderIDs))

	for i := 0; i < s.workerCount; i++ {
		group.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case orderID, ok := <-orderIDsChan:
					if !ok {
						return nil
					}

					orderIDWithExpiryAndStatus, err := s.getOrderIDWithExpiryAndStatus(ctx, orderID, currTime)
					if err != nil {
						return err
					}
					resChan <- orderIDWithExpiryAndStatus
				}
			}
		})
	}

	for _, orderID := range orderIDs {
		orderIDsChan <- orderID
	}
	close(orderIDsChan)

	if err := group.Wait(); err != nil {
		defer close(resChan)
		return nil, err
	}
	close(resChan)

	res := make([]OrderIDWithExpiryAndStatus, 0, len(resChan))
	for orderIDWithExpiryAndStatus := range resChan {
		res = append(res, orderIDWithExpiryAndStatus)
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

// GetReturnsList method is generator, it uses chan to return results from db
func (s *Service) GetReturnsList(ctx context.Context, offset, limit int) (<-chan ReturnOrderAndCustomer, error) {
	if offset < 0 {
		return nil, ErrorInvalidOffsetValue
	}

	if limit <= 0 {
		return nil, ErrorInvalidLimitValue
	}

	returnOrderIDs, err := s.orderStorage.GetOrderIDsWhereStatus(ctx, models.StatusReturn, uint(offset), uint(limit))
	if err != nil {
		return nil, err
	}

	// Worker Pool
	group, ctx := errgroup.WithContext(ctx)
	orderIDsChan := make(chan models.IDType, len(returnOrderIDs))
	orderAndCustomerChan := make(chan ReturnOrderAndCustomer, len(returnOrderIDs))

	for i := 0; i < s.workerCount; i++ {
		group.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case orderID, ok := <-orderIDsChan:
					if !ok {
						return nil
					}

					customerID, err := s.orderStorage.GetOrderCustomerID(ctx, orderID)
					if err != nil {
						return err
					}

					orderAndCustomerChan <- ReturnOrderAndCustomer{
						OrderID:    orderID,
						CustomerID: customerID,
					}
				}
			}
		})
	}

	for _, orderID := range returnOrderIDs {
		orderIDsChan <- orderID
	}
	close(orderIDsChan)

	defer close(orderAndCustomerChan)
	if err := group.Wait(); err != nil {
		return nil, err
	}

	return orderAndCustomerChan, nil
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

func (s *Service) ChangeWorkerCount(workerCount int) error {
	if workerCount <= 0 || workerCount > runtime.GOMAXPROCS(0) {
		return ErrorInvalidWorkerCount
	}

	s.workerCount = workerCount

	return nil
}

func (s *Service) GetWorkerCount() int {
	return s.workerCount
}
