package startProducer

import (
	"context"
	"encoding/json"
	"log"
	"producer/models"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// SendTestMessage launches cmd producer and sends message to consumer
func SendTestMessage() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "test123",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

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

func createValidData() *models.CombinedData {
	orderUUID := uuid.New().String()
	deliveryUUID := uuid.New().String()

	orderTime := time.Now()

	track := "WBILMTESTTRACK"
	entry := "WBIL"
	locale := "en"
	internalSignature := ""
	customer := "test"
	deliveryService := "meest"
	shardKey := "9"
	smID := 99
	oofShard := "1"

	name := "Test Testov"
	phone := "+9720000000"
	zip := "2639809"
	city := "Kiryat Mozkin"
	address := "Ploshad Mira 15"
	region := "Kraiot"
	email := "test@gmail.com"

	requestId := ""
	currency := "USD"
	provider := "wbpay"
	amount := 1817
	paymentDT := 1637907727
	bank := "alpha"
	deliveryCost := 1500
	goodsTotal := 317
	customFee := 0

	chrtID := 9934930
	price := 453
	rid := "ab4219087a764ae0btest"
	name2 := "Mascaras"
	sale := 30
	size := "0"
	totalPrice := 317
	nmId := 2389212
	brand := "Vivienne Sabo"
	status := 202

	data := models.CombinedData{
		Order: models.Order{
			orderUUID,
			&track,
			&entry,
			&deliveryUUID,
			&locale,
			&internalSignature,
			&customer,
			&deliveryService,
			&shardKey,
			&smID,
			&orderTime,
			&oofShard,
		},
		Delivery: models.Delivery{
			&deliveryUUID,
			&name,
			&phone,
			&zip,
			&city,
			&address,
			&region,
			&email,
		},
		Payment: models.Payment{
			&orderUUID,
			&requestId,
			&currency,
			&provider,
			&amount,
			&paymentDT,
			&bank,
			&deliveryCost,
			&goodsTotal,
			&customFee,
		},
		Items: []models.Item{models.Item{
			&chrtID,
			&track,
			&price,
			&rid,
			&name2,
			&sale,
			&size,
			&totalPrice,
			&nmId,
			&brand,
			&status,
		}},
	}

	return &data
}
