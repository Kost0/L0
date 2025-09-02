package startProducer

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"producer/models"
	"sync"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/segmentio/kafka-go"
)

// StartProducer launches kafka producer
func StartProducer() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "test1",
		Balancer: &kafka.LeastBytes{},
	}
	defer func() {
		if err := writer.Close(); err != nil {
			log.Println(err)
		}
	}()

	rand.Seed(time.Now().UnixNano())

	wg := &sync.WaitGroup{}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down Kafka server...")
				return
			default:
				sendTestMessage(writer)
				time.Sleep(time.Second)
				sendWrongMessage(writer)
				time.Sleep(time.Second)
			}
		}
	}()

	wg.Wait()
}

func sendTestMessage(writer *kafka.Writer) {
	data := createValidData()

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

	print(data.Order.OrderUID)
}

func sendWrongMessage(writer *kafka.Writer) {
	data := createWrongData()

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
}

func createWrongData() *models.CombinedData {
	return &models.CombinedData{}
}

func createValidData() *models.CombinedData {
	order := &models.Order{}
	pay := &models.Payment{}
	deliv := &models.Delivery{}
	item := &models.Item{}

	err := gofakeit.Struct(order)
	if err != nil {
		log.Fatal()
	}
	err = gofakeit.Struct(pay)
	if err != nil {
		log.Fatal()
	}
	err = gofakeit.Struct(deliv)
	if err != nil {
		log.Fatal()
	}
	err = gofakeit.Struct(item)
	if err != nil {
		log.Fatal()
	}

	order.DeliveryID = deliv.ID
	pay.Transaction = &order.OrderUID
	order.TrackNumber = item.TrackNumber

	data := &models.CombinedData{
		*order,
		*pay,
		*deliv,
		[]models.Item{*item},
	}

	return data
}
