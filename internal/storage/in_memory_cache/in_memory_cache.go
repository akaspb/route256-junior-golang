package in_memory_cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/in_memory_cache"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/in_memory_cache/ttl_cache"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
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

// NewInMemoryCache создаст указатель на структуру InMemoryCache, которая
// имплементирует интерфейс storage.Facade
//
// Задача этой структуры - создавать in-memory кеш, чтобы снизить нагрузку на БД
func NewInMemoryCache(repository storage.Facade, validTime time.Duration) *InMemoryCache {
	return &InMemoryCache{
		repo:  repository,
		cache: ttl_cache.NewTTLCache[string, any](validTime),
		mu:    sync.Mutex{},
	}
}

// CreateOrder добавляет информацию о заказе в БД и in-memory кеш.
//
// Скорость добавления в in-memory кеш O(1)
func (c *InMemoryCache) CreateOrder(ctx context.Context, order models.Order) error {
	c.setOrder(order)
	c.incCounter(order.CustomerID, 1)
	return c.repo.CreateOrder(ctx, order)
}

// GetOrder ищет информацию о заказе в in-memory кеше. Если информация есть и валидна, то обращения к БД не происходит
//
// Скорость поиска в in-memory кеше O(1)
func (c *InMemoryCache) GetOrder(ctx context.Context, orderID models.IDType) (models.Order, error) {
	order, ok := c.getOrder(orderID)
	if !ok {
		return c.repo.GetOrder(ctx, orderID)
	}

	return order, nil
}

// DeleteOrder удаляет информацию о заказе из in-memory кеша и из БД
//
// Скорость удаления в in-memory кеше O(1)
func (c *InMemoryCache) DeleteOrder(ctx context.Context, orderID models.IDType) error {
	order, ok := c.getOrder(orderID)
	if ok {
		c.decCounter(order.CustomerID, 1)
		c.deleteOrder(orderID)
	}
	return c.repo.DeleteOrder(ctx, orderID)
}

// ChangeOrderStatus изменяет значение статуса заказа в in-memory кеше и в БД
//
// Скорость выполнения в in-memory кеше O(1)
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

// GetCustomerOrderIDsWithStatus получает ID заказов по ID владельца с заданным значением статуса заказа
// Если информация валидна (т.е. не устарела), то обращения к БД не будет
//
// Скорость выполнения в in-memory кеше O({количество заказов владельца со статусом из кеша})
func (c *InMemoryCache) GetCustomerOrderIDsWithStatus(ctx context.Context, customerID models.IDType, statusVal models.StatusVal) ([]models.IDType, error) {
	slc := c.getAllCustomerOrdersWithStatus(customerID, statusVal)
	if slc != nil {
		return slc, nil
	}

	return c.repo.GetCustomerOrderIDsWithStatus(ctx, customerID, statusVal)
}

// GetNCustomerOrderIDsWithStatus получает ID заказов по ID владельца с заданным значением статуса заказа с ограничением
// максимальной длины результата n.
// Если информация валидна (т.е. не устарела), то обращения к БД не будет
//
// Скорость выполнения в in-memory кеше O(max({количество заказов владельца со статусом из кеша}, n))
func (c *InMemoryCache) GetNCustomerOrderIDsWithStatus(ctx context.Context, customerID models.IDType, statusVal models.StatusVal, n uint) ([]models.IDType, error) {
	slc := c.getMaxNCustomerOrdersWithStatus(customerID, statusVal, int(n))
	if slc != nil {
		return slc, nil
	}

	return c.repo.GetNCustomerOrderIDsWithStatus(ctx, customerID, statusVal, n)
}

// GetOrderStatus значение статуса заказа по его ID
// максимальной длины результата n.
// Если информация валидна (т.е. не устарела), то обращения к БД не будет
//
// Скорость выполнения в in-memory кеше O(1)
func (c *InMemoryCache) GetOrderStatus(ctx context.Context, orderID models.IDType) (models.Status, error) {
	order, ok := c.getOrder(orderID)
	if !ok {
		return c.repo.GetOrderStatus(ctx, orderID)
	}

	return order.Status, nil
}

// GetOrderIDsWhereStatus получает заданное число ID заказов с заданным значением статуса заказа
//
// Метод работает следующим образом:
//   - получаем список владельцев
//   - просматриваем такое количество id владельцев, чтобы offset стал меньше 0 на этом владельце
//   - просматриваем заказы последнего рассматриваемого владельца, чтобы начать добавлять с того id заказа,
//     согласно offset
//   - если в какой-то момент получилось набрать limit число заказов, то возвращаем результат,
//
// иначе смотрим в БД
//
// Скорость выполнения в in-memory кеше O(offset + limit)
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

// GetOrderCustomerID возвращает владельца заказа. Метод ищет информацию о заказе в in-memory кеше.
// Если информация есть и валидна, то обращения к БД не происходит
//
// Скорость поиска в in-memory кеше O(1)
func (c *InMemoryCache) GetOrderCustomerID(ctx context.Context, orderID models.IDType) (models.IDType, error) {
	order, ok := c.getOrder(orderID)
	if !ok {
		return c.repo.GetOrderCustomerID(ctx, orderID)
	}

	return order.CustomerID, nil
}

// getOrder обращается к in-memory кешу для получения информации об заказе
func (c *InMemoryCache) getOrder(orderId models.IDType) (models.Order, bool) {
	now := time.Now()
	orderCacheKey := fmt.Sprintf("%s:%v", orderPrefix, orderId)
	res, ok := c.cache.Get(orderCacheKey, now)
	if !ok {
		return models.Order{}, false
	}
	return res.(models.Order), true
}

// setOrder обращается к in-memory кешу для добавления (перезаписи) информации об заказе
func (c *InMemoryCache) setOrder(order models.Order) {
	now := time.Now()
	orderCacheKey := fmt.Sprintf("%s:%v", orderPrefix, order.ID)
	c.cache.Set(orderCacheKey, order, now)
}

// deleteOrder обращается к in-memory кешу для удаления информации об заказе
func (c *InMemoryCache) deleteOrder(orderId models.IDType) bool {
	orderCacheKey := fmt.Sprintf("%s:%v", orderPrefix, orderId)
	return c.cache.Delete(orderCacheKey)
}

// countCustomerOrdersWithStatus обращается к in-memory кешу для получения числа заказов у владельца с
// заданным значением статуса заказа
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

// addCustomerOrderWithStatus обращается к in-memory кешу для добавления владельцу id заказа в список с заданным значением статуса заказа
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

// addCustomerOrderWithStatus обращается к in-memory кешу для получения всех заказов владельца с заданным значением статуса
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

// addCustomerOrderWithStatus обращается к in-memory кешу для получения ограниченного числом n заказов владельца с заданным значением статуса
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

// InMemoryCache проводит инвалидацию информации о заказах владельца с заданным значением статуса
//
// Метод используется в методе ChangeOrderStatus, т.к. было бы затратно по времени искать и удалять
// ID заказа, у которого поменялось значение статуса
func (c *InMemoryCache) invalidateCustomerOrdersWithStatus(customerID models.IDType, statusVal models.StatusVal) {
	delta := c.countCustomerOrdersWithStatus(customerID, statusVal)
	c.decCounter(customerID, delta)

	customerIDStatusValKey := fmt.Sprintf("%s:%v:%v", customerPrefix, customerID, statusVal)
	c.cache.Delete(customerIDStatusValKey)
}

// incCounter увеличивает значение в счётчике числа заказов владельца, про которые есть информация в in-memory кеше
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

// incCounter уменьшает значение в счётчике числа заказов владельца, про которые есть информация в in-memory кеше
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

// getCustomerIDs получает список владельцев, у которых не 0 заказов в in-memory кеше
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
