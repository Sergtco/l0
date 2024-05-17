package cache

import (
	"fmt"
	"service/data/database"
	"service/models"
	"sync"
)

var (
	NotExistError error = fmt.Errorf("No item exists.")
)

type Cache struct {
	data     map[string]*models.Order
	lastUsed []string
	maxSize  int
	dataBase *database.Database
	mu       sync.RWMutex
}

func (c *Cache) Get(key string) (*models.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.data[key]; ok {
		return val, nil
	} else if val, err := c.dataBase.GetOrder(key); err == nil {
		go func() {
			c.put(key, val)
		}()
		return val, nil
	}
	return nil, NotExistError
}

func (c *Cache) put(key string, value *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.maxSize >= len(c.lastUsed) && len(c.lastUsed) != 0 {
		lastKey := c.lastUsed[0]
		c.lastUsed = c.lastUsed[1:]
		delete(c.data, lastKey)
	}
	if _, ok := c.data[key]; !ok {
		c.data[key] = value
		c.lastUsed = append(c.lastUsed, key)
	}
}

func (c *Cache) Store(key string, value *models.Order) error {
	err := c.dataBase.InsertOrder(value)
	if err != nil {
		return err
	}
	go func() {
		c.put(key, value)
	}()
	return nil
}

func NewCache(maxSize int) *Cache {
	c := &Cache{
		data:     map[string]*models.Order{},
		lastUsed: make([]string, 0, maxSize),
		maxSize:  maxSize,
		dataBase: database.New(),
		mu:       sync.RWMutex{},
	}
	c.restoreData()
	return c
}
func (c *Cache) restoreData() {
	c.dataBase.GetTopOrders(c.maxSize)
}
