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

type OrderCache struct {
	data sync.Map
	ttl  time.Duration
}

func NewOrderCache(ttl time.Duration) *OrderCache {
	return &OrderCache{ttl: ttl}
}

func (c *OrderCache) Set(orderID string, data *models.CombinedData) {
	c.data.Store(orderID, data)
	time.AfterFunc(c.ttl, func() {
		c.data.Delete(orderID)
	})
}

func (c *OrderCache) Get(orderID string) (models.CombinedData, bool) {
	v, ok := c.data.Load(orderID)
	if !ok {
		return models.CombinedData{}, false
	}
	return v.(models.CombinedData), true
}

func (c *OrderCache) WarmUpCache(db *sql.DB, ctx context.Context) error {
	data, err := getRecentOrders(db, ctx)
	if err != nil {
		return err
	}

	for _, order := range data {
		c.Set(order.Order.OrderUID, order)
	}

	log.Printf("Warmed up cache with %d orders", len(data))

	return nil
}

func getRecentOrders(db *sql.DB, ctx context.Context) ([]*models.CombinedData, error) {
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

	for rows.Next() {
		id := ""
		err = rows.Scan(id)
		if err != nil {
			return nil, err
		}
		data, err := repository.SelectOrder(db, id)
		if err != nil {
			return nil, err
		}
		allData = append(allData, data)
	}

	return allData, nil
}
