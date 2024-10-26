package in_memory_cache

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/in_memory_cache"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/in_memory_cache/ttl_cache"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"sync"
	"time"
)

const (
	orderPrefix    = "order"
	customerPrefix = "customer"
	orderCounter   = "counter"
)

type InMemoryCache struct {
	repo  storage.Facade
	cache *in_memory_cache.InMemoryCache[string, time.Time, any]
	mu    sync.Mutex
}

func NewInMemoryCache(repository storage.Facade, validTime time.Duration) *InMemoryCache {
	return &InMemoryCache{
		repo:  repository,
		cache: ttl_cache.NewTTLCache[string, any](validTime),
		mu:    sync.Mutex{},
	}
}

func (c *InMemoryCache) CreateOrder(ctx context.Context, order models.Order) error {
	c.setOrder(order)
	c.incCounter(order.CustomerID, 1)
	return c.repo.CreateOrder(ctx, order)
}

func (c *InMemoryCache) GetOrder(ctx context.Context, orderID models.IDType) (models.Order, error) {
	order, ok := c.getOrder(orderID)
	if !ok {
		return c.repo.GetOrder(ctx, orderID)
	}

	return order, nil
}

func (c *InMemoryCache) DeleteOrder(ctx context.Context, orderID models.IDType) error {
	order, ok := c.getOrder(orderID)
	if ok {
		c.decCounter(order.CustomerID, 1)
		c.deleteOrder(orderID)
	}
	return c.repo.DeleteOrder(ctx, orderID)
}

func (c *InMemoryCache) ChangeOrderStatus(ctx context.Context, orderID models.IDType, val models.Status) (err error) {
	err = c.repo.ChangeOrderStatus(ctx, orderID, val)

	order, ok := c.getOrder(orderID)
	if !ok {
		return
	}

	c.invalidateCustomerOrdersWithStatus(order.CustomerID, order.Status.Value)

	order.Status = val
	c.setOrder(order)

	return
}

// status ChangedAt value is important here
func (c *InMemoryCache) GetCustomerOrderIDsWithStatus(ctx context.Context, customerID models.IDType, statusVal models.StatusVal) ([]models.IDType, error) {
	slc := c.getAllCustomerOrdersWithStatus(customerID, statusVal)
	if slc != nil {
		return slc, nil
	}

	return c.repo.GetCustomerOrderIDsWithStatus(ctx, customerID, statusVal)
}

// status ChangedAt value is important here
func (c *InMemoryCache) GetNCustomerOrderIDsWithStatus(ctx context.Context, customerID models.IDType, statusVal models.StatusVal, n uint) ([]models.IDType, error) {
	slc := c.getMaxNCustomerOrdersWithStatus(customerID, statusVal, int(n))
	if slc != nil {
		return slc, nil
	}

	return c.repo.GetNCustomerOrderIDsWithStatus(ctx, customerID, statusVal, n)
}

func (c *InMemoryCache) GetOrderStatus(ctx context.Context, orderID models.IDType) (models.Status, error) {
	order, ok := c.getOrder(orderID)
	if !ok {
		return c.repo.GetOrderStatus(ctx, orderID)
	}

	return order.Status, nil
}

func (c *InMemoryCache) GetOrderIDsWhereStatus(ctx context.Context, statusVal models.StatusVal, offset, limit uint) ([]models.IDType, error) {
	customerIDs := c.getCustomerIDs()
	res := make([]models.IDType, 0, limit)

	i := 0
	for ; i < len(customerIDs); i++ {
		customerID := customerIDs[i]
		count := c.countCustomerOrdersWithStatus(customerID, statusVal)
		if int(offset)-count < 0 {
			break
		}
		offset -= uint(count)
	}

	for ; i < len(customerIDs); i++ {
		customerID := customerIDs[i]
		orderIDs := c.getAllCustomerOrdersWithStatus(customerID, statusVal)
		for _, orderID := range orderIDs {
			if offset != 0 {
				offset--
			} else {
				res = append(res, orderID)
				if len(res) == int(limit) {
					return res, nil
				}
			}
		}
	}

	return c.repo.GetOrderIDsWhereStatus(ctx, statusVal, offset, limit)
}

func (c *InMemoryCache) GetOrderCustomerID(ctx context.Context, orderID models.IDType) (models.IDType, error) {
	order, ok := c.getOrder(orderID)
	if !ok {
		return c.repo.GetOrderCustomerID(ctx, orderID)
	}

	return order.CustomerID, nil
}

func (c *InMemoryCache) getOrder(orderId models.IDType) (models.Order, bool) {
	now := time.Now()
	orderCacheKey := fmt.Sprintf("%s:%v", orderPrefix, orderId)
	res, ok := c.cache.Get(orderCacheKey, now)
	if !ok {
		return models.Order{}, false
	}
	return res.(models.Order), true
}

func (c *InMemoryCache) setOrder(order models.Order) {
	now := time.Now()
	orderCacheKey := fmt.Sprintf("%s:%v", orderPrefix, order.ID)
	c.cache.Set(orderCacheKey, order, now)
}

func (c *InMemoryCache) deleteOrder(orderId models.IDType) bool {
	orderCacheKey := fmt.Sprintf("%s:%v", orderPrefix, orderId)
	return c.cache.Delete(orderCacheKey)
}

func (c *InMemoryCache) countCustomerOrdersWithStatus(customerID models.IDType, statusVal models.StatusVal) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	customerIDStatusValKey := fmt.Sprintf("%s:%v:%v", customerPrefix, customerID, statusVal)

	idsAny, ok := c.cache.Get(customerIDStatusValKey, now)
	if !ok {
		return 0
	}

	return len(idsAny.([]models.IDType))
}

func (c *InMemoryCache) addCustomerOrderWithStatus(orderID, customerID models.IDType, statusVal models.StatusVal) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	customerIDStatusValKey := fmt.Sprintf("%s:%v:%v", customerPrefix, customerID, statusVal)

	idsAny, ok := c.cache.Get(customerIDStatusValKey, now)
	if !ok {
		c.cache.Set(customerIDStatusValKey, []models.IDType{orderID}, now)
		return
	}

	ids := idsAny.([]models.IDType)
	ids = append(ids, orderID)
	c.cache.Set(customerIDStatusValKey, ids, now)
}

func copyLastNDesc[T any](slc []T, n int) []T {
	res := make([]T, n)
	for i := 0; i < n; i++ {
		res[i] = slc[n-1-i]
	}

	return res
}

func (c *InMemoryCache) getAllCustomerOrdersWithStatus(customerID models.IDType, statusVal models.StatusVal) []models.IDType {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	customerIDStatusValKey := fmt.Sprintf("%s:%v:%v", customerPrefix, customerID, statusVal)

	idsAny, ok := c.cache.Get(customerIDStatusValKey, now)
	if !ok {
		return nil
	}

	ids := idsAny.([]models.IDType)
	return copyLastNDesc(ids, len(ids))
}

func (c *InMemoryCache) getMaxNCustomerOrdersWithStatus(customerID models.IDType, statusVal models.StatusVal, n int) []models.IDType {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	customerIDStatusValKey := fmt.Sprintf("%s:%v:%v", customerPrefix, customerID, statusVal)

	idsAny, ok := c.cache.Get(customerIDStatusValKey, now)
	if !ok {
		return nil
	}

	ids := idsAny.([]models.IDType)
	return copyLastNDesc(ids, max(len(ids), n))
}

func (c *InMemoryCache) invalidateCustomerOrdersWithStatus(customerID models.IDType, statusVal models.StatusVal) {
	delta := c.countCustomerOrdersWithStatus(customerID, statusVal)
	c.decCounter(customerID, delta)

	customerIDStatusValKey := fmt.Sprintf("%s:%v:%v", customerPrefix, customerID, statusVal)
	c.cache.Delete(customerIDStatusValKey)
}

func (c *InMemoryCache) incCounter(customerID models.IDType, delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	counterAny, ok := c.cache.Get(orderCounter, now)
	counter := make(map[models.IDType]int, 1)
	if ok {
		counter = counterAny.(map[models.IDType]int)
	}

	counter[customerID] += delta

	c.cache.Set(orderCounter, counter, now)
}

func (c *InMemoryCache) decCounter(customerID models.IDType, delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	counterAny, ok := c.cache.Get(orderCounter, now)
	if !ok {
		return
	}

	counter := counterAny.(map[models.IDType]int)
	counter[customerID] -= delta
	if counter[customerID] <= 0 {
		delete(counter, customerID)
	}

	c.cache.Set(orderCounter, counter, now)
}

func (c *InMemoryCache) getCustomerIDs() []models.IDType {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	counterAny, ok := c.cache.Get(orderCounter, now)
	if !ok {
		return nil
	}

	counter := counterAny.(map[models.IDType]int)
	res := make([]models.IDType, 0, len(counter))
	for id, _ := range counter {
		res = append(res, id)
	}

	return res
}
