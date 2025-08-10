package repository

import (
	"database/sql"
	"encoding/json"
	"github.com/Kost0/L0/internal/models"
	"github.com/segmentio/kafka-go"
)

type CombinedData struct {
	order    models.Order
	payment  models.Payment
	delivery models.Delivery
	items    []models.Item
}

func InsertOrder(db *sql.DB, msg *kafka.Message) error {
	var data CombinedData
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		return err
	}

	return nil
}
