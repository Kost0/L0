package kafka

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Kost0/L0/internal/repository"
	"github.com/segmentio/kafka-go"
)

type DLQHandler struct {
	mainReader *kafka.Reader
	dlqWriter  *kafka.Writer
	maxRetries int
}

func NewDLQHandler(brokers []string, mainTopic, dlqTopic string) *DLQHandler {
	return &DLQHandler{
		mainReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   mainTopic,
			GroupID: "main",
		}),
		dlqWriter: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    dlqTopic,
			Balancer: &kafka.Hash{},
		},
		maxRetries: 3,
	}
}

func (h *DLQHandler) ProcessWithRetry(ctx context.Context, repo repository.OrderRepository) error {
	for {
		msg, err := h.mainReader.ReadMessage(context.Background())
		if err != nil {
			if ctx.Err() != nil || errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return err
			}
			log.Printf("Error reading message: %s\n", err)
			continue
		}

		if err = h.processWithRetry(ctx, repo, &msg); err != nil {
			log.Printf("Final processing error: %v", err)
		}
	}
}

func (h *DLQHandler) processWithRetry(ctx context.Context, repo repository.OrderRepository, msg *kafka.Message) error {
	for attempt := 1; attempt <= h.maxRetries; attempt++ {
		err := processMessage(ctx, repo, msg)
		if err == nil {
			return nil
		}

		log.Printf("Attempt #%d: %v", attempt, err)

		if attempt == h.maxRetries {
			return h.sendToDLQ(msg, err)
		}

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return nil
}

func (h *DLQHandler) sendToDLQ(msg *kafka.Message, err error) error {
	dlqMsg := kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
		Headers: append(msg.Headers,
			kafka.Header{Key: "original_topic", Value: []byte(msg.Topic)},
			kafka.Header{Key: "error", Value: []byte(err.Error())},
			kafka.Header{Key: "timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		),
		Time: msg.Time,
	}

	return h.dlqWriter.WriteMessages(context.Background(), dlqMsg)
}
