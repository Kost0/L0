// Package cache provides work with cache
//
// Includes:
//   - creating a cache
//   - getting from cache
//   - filling in data at the start of the program
package cache

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/Kost0/L0/internal/models"
	"github.com/Kost0/L0/internal/repository"
)

// Cache defines interface for working with cache
type Cache interface {
	Set(orderID string, data *models.CombinedData)
	Get(orderID string) (*models.CombinedData, bool)
	WarmUpCache(db *sql.DB, repo repository.OrderRepository, ctx context.Context) error
}

// OrderCache contains the location and time of storage of the cache
type OrderCache struct {
	data sync.Map
	ttl  time.Duration
}

// NewOrderCache create new OrderCache
// Accepts:
//   - ttl: time to live
//
// Returns:
//   - *OrderCache
func NewOrderCache(ttl time.Duration) *OrderCache {
	return &OrderCache{ttl: ttl}
}

// Set save data in cache
// Accepts:
//   - orderID: id of order
//   - data: all data about order
func (c *OrderCache) Set(orderID string, data *models.CombinedData) {
	c.data.Store(orderID, data)
	time.AfterFunc(c.ttl, func() {
		c.data.Delete(orderID)
	})
}

// Get receive data from cache
// Accepts:
//   - orderID: id of order
//
// Returns:
//   - all data about order
//   - did it work
func (c *OrderCache) Get(orderID string) (*models.CombinedData, bool) {
	v, ok := c.data.Load(orderID)
	if !ok {
		return &models.CombinedData{}, false
	}
	return v.(*models.CombinedData), true
}

// WarmUpCache
// Accepts:
//   - db: database
//   - repo: repository
//   - cts: context
//
// Returns:
//   - error if something wrong
func (c *OrderCache) WarmUpCache(db *sql.DB, repo repository.OrderRepository, ctx context.Context) error {
	data, err := getRecentOrders(db, repo, ctx)
	if err != nil {
		return err
	}

	for _, order := range data {
		c.Set(order.Order.OrderUID, order)
	}

	log.Printf("Warmed up cache with %d orders", len(data))

	return nil
}

func getRecentOrders(db *sql.DB, repo repository.OrderRepository, ctx context.Context) ([]*models.CombinedData, error) {
	query := `
SELECT order_uid FROM orders
WHERE date_created >= NOW() - INTERVAL '7 days'
`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allData := []*models.CombinedData{}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	for rows.Next() {
		id := ""
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		data, err := repo.SelectWithRetry(ctx, id)
		if err != nil {
			return nil, err
		}
		allData = append(allData, data)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allData, nil
}
