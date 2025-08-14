package repository

import (
	"github.com/Kost0/L0/internal/models"
)

type OrderRepository interface {
	SelectOrder(orderID string) (*models.CombinedData, error)
}
