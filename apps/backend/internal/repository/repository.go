// Package repository provides work with database
//
// Includes:
//   - connecting to database
//   - migrations
//   - operators
package repository

import (
	"context"

	"github.com/Kost0/L0/internal/models"
)

// OrderRepository defines interface for working with database
type OrderRepository interface {
	SelectOrder(orderID string) (*models.CombinedData, error)
	SelectWithRetry(ctx context.Context, orderUID string) (*models.CombinedData, error)
	InsertOrder(data *models.CombinedData) error
	InsertWithRetry(ctx context.Context, data *models.CombinedData) error
}
