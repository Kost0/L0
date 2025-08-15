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

	orderTime, err := time.Parse(time.RFC3339, "2025-08-11T06:22:19Z")
	if err != nil {
		log.Fatal(err)
	}

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

	transaction := "b563feb7b2b84b6test"
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
	trackNum := "WBILMTESTTRACK"
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
			&paymentUUID,
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
			&paymentUUID,
			&transaction,
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
			&ItemUUID,
			&orderUUID,
			&chrtID,
			&trackNum,
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
