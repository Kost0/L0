package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Kost0/L0/internal/models"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func SendTestMessage() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	orderUUID := uuid.New().String()
	deliveryUUID := uuid.New().String()
	paymentUUID := uuid.New().String()
	ItemUUID := uuid.New().String()

	orderTime, err := time.Parse(time.RFC3339, "2021-11-26T06:22:19Z")
	if err != nil {
		log.Fatal(err)
	}

	data := models.CombinedData{
		Order: models.Order{
			orderUUID,
			"WBILMTESTTRACK",
			"WBIL",
			deliveryUUID,
			paymentUUID,
			"en",
			"",
			"test",
			"meest",
			"9",
			99,
			orderTime,
			"1",
		},
		Delivery: models.Delivery{
			deliveryUUID,
			"Test Testov",
			"+9720000000",
			"2639809",
			"Kiryat Mozkin",
			"Ploshad Mira 15",
			"Kraiot",
			"test@gmail.com",
		},
		Payment: models.Payment{
			paymentUUID,
			"b563feb7b2b84b6test",
			"",
			"USD",
			"wbpay",
			1817,
			1637907727,
			"alpha",
			1500,
			317,
			0,
		},
		Items: []models.Item{models.Item{
			ItemUUID,
			orderUUID,
			9934930,
			"WBILMTESTTRACK",
			453,
			"ab4219087a764ae0btest",
			"Mascaras",
			30,
			"0",
			317,
			2389212,
			"Vivienne Sabo",
			202,
		}},
	}

	buf, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("test"),
		Value: buf,
		Time:  time.Now(),
	})

	if err != nil {
		log.Printf("Failed to write messages: %s", err)
	} else {
		log.Printf("Successfully sent message")
	}
	print(orderUUID)
}
