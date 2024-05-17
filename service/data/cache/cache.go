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

// Base cache with queue, uses mutex to prevent race condition/
type Cache struct {
	data     map[string]*models.Order
	lastUsed []string
	maxSize  int
	dataBase *database.Database
	mu       sync.RWMutex
}

// Returns NotExistError if key is not valid
func (c *Cache) Get(key string) (*models.Order, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.data[key]; ok {
		return val, nil
	} else if val, err := c.dataBase.GetOrder(key); err == nil {
		c.put(val)
		return val, nil
	}
	return nil, NotExistError
}

func (c *Cache) put(value *models.Order) {
	if c.maxSize <= len(c.lastUsed) && len(c.lastUsed) != 0 {
		lastKey := c.lastUsed[0]
		c.lastUsed = c.lastUsed[1:]
		delete(c.data, lastKey)
	}
	if _, ok := c.data[value.Id]; !ok {
		c.data[value.Id] = value
		c.lastUsed = append(c.lastUsed, value.Id)
	}
}

func (c *Cache) Store(value *models.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.dataBase.InsertOrder(value)
	if err != nil {
		return err
	}
	c.put(value)
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
	fmt.Println(len(c.data))
	return c
}

func (c *Cache) restoreData() error {
	orders, err := c.dataBase.GetTopOrders(c.maxSize)
	if err != nil {
		return err
	}
	for _, order := range orders {
		c.Store(order)
	}
	return nil
}
