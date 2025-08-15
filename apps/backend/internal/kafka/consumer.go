package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/mail"
	"time"

	"github.com/Kost0/L0/internal/models"
	"github.com/Kost0/L0/internal/repository"
	"github.com/go-playground/validator"
	"github.com/segmentio/kafka-go"
)

func StartKafka(ctx context.Context, db *sql.DB) {
	brokerAddress := "kafka:19092"
	topic := "orders"
	groupID := "myGroup"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          []string{brokerAddress},
		Topic:            topic,
		GroupID:          groupID,
		MinBytes:         10e3,
		MaxBytes:         10e6,
		MaxWait:          1 * time.Second,
		RebalanceTimeout: 20 * time.Second,
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down Kafka server...")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			log.Print("Reading message...")
			if err != nil {
				if ctx.Err() != nil || ctx.Err() == context.DeadlineExceeded {
					return
				}
				log.Printf("Error reading message: %s\n", err)
				continue
			}

			var data models.CombinedData
			err = json.Unmarshal(msg.Value, &data)
			if err != nil {
				log.Printf("Error unmarshalling message: %s\n", err)
			}

			if err = validateData(&data); err != nil {
				log.Printf("Error validating data: %s\n", err)
			}

			err = repository.InsertWithRetry(ctx, db, &data)
			if err != nil {
				log.Printf("Error inserting order: %s\n", err)
			}

			log.Printf("Message on %s: %s\n", msg.Topic, string(msg.Value))
		}
	}
}

func validateData(data *models.CombinedData) error {
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		return err
	}

	if data.Order.DateCreated.After(time.Now()) {
		return errors.New("Order date created is in the future")
	}

	if _, err := mail.ParseAddress(*data.Delivery.Email); err != nil {
		return err
	}

	for _, item := range data.Items {
		if *item.Price < 0 {
			return errors.New("Item price is negative")
		}
	}

	return nil
}
