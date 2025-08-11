package kafka

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Kost0/L0/internal/repository"
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
			if err != nil {
				if ctx.Err() != nil || ctx.Err() == context.DeadlineExceeded {
					return
				}
				log.Printf("Error reading message: %s\n", err)
				continue
			}
			err = repository.InsertOrder(db, &msg)
			if err != nil {
				log.Printf("Error inserting order: %s\n", err)
			}
			log.Printf("Message on %s: %s\n", msg.Topic, string(msg.Value))
		}
	}
}
