// Package cmd provides work with cmd
package kafka

import (
	"context"
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

// StartKafka launches cmd consumer to process messages
// Accepts:
//   - ctx: context
//   - db: database
func StartKafka(ctx context.Context, repo repository.OrderRepository) {
	brokerAddress := "kafka:9092"
	topic := "test123"
	groupID := "myOrdersGroup-12345"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          []string{brokerAddress},
		Topic:            topic,
		GroupID:          groupID,
		MinBytes:         10e3,
		MaxBytes:         10e6,
		MaxWait:          1 * time.Second,
		RebalanceTimeout: 20 * time.Second,
		StartOffset:      kafka.FirstOffset,
		CommitInterval:   0,
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down Kafka server...")
			return
		default:
			msg, err := reader.ReadMessage(ctx)

			if err != nil {
				if ctx.Err() != nil || errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return
				}
				log.Printf("Error reading message: %s\n", err)
				continue
			}
			err = processMessage(ctx, repo, &msg)
			if err != nil {
				log.Printf("Error processing message: %s\n", err)
			}
		}
	}
}

func processMessage(ctx context.Context, repo repository.OrderRepository, msg *kafka.Message) error {
	var data models.CombinedData
	log.Printf("Received message: %s\n", string(msg.Value))
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		log.Printf("Error unmarshalling message: %s\n", err)
		return err
	}

	if err = validateData(&data); err != nil {
		log.Printf("Error validating data: %s\n", err)
		return err
	}

	err = repo.InsertWithRetry(ctx, &data)
	if err != nil {
		log.Printf("Error inserting order: %s\n", err)
		return err
	}

	log.Printf("Message on %s: %s\n", msg.Topic, string(msg.Value))
	return nil
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
