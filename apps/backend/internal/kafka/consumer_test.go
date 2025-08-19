package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Kost0/L0/internal/models"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestValidateData_ValidData(t *testing.T) {
	data := createValidData()

	err := validateData(data)
	assert.NoError(t, err)
}

func TestValidateData_FutureDate(t *testing.T) {
	data := createValidData()
	futureTime := data.Order.DateCreated.Add(2 * time.Hour)
	data.Order.DateCreated = &futureTime

	err := validateData(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "future")
}

func TestValidateData_InvalidEmail(t *testing.T) {
	data := createValidData()
	wrongEmail := ""

	data.Delivery.Email = &wrongEmail

	err := validateData(data)
	assert.Error(t, err)
}

func TestValidateData_NegativePrice(t *testing.T) {
	data := createValidData()
	wrongPrice := -100

	data.Items[0].Price = &wrongPrice

	err := validateData(data)
	assert.Error(t, err)
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) SelectOrder(orderID string) (*models.CombinedData, error) {
	args := m.Called(orderID)
	return args.Get(0).(*models.CombinedData), args.Error(1)
}

func (m *MockOrderRepository) SelectWithRetry(ctx context.Context, orderID string) (*models.CombinedData, error) {
	args := m.Called(orderID)
	return args.Get(0).(*models.CombinedData), args.Error(1)
}

func (m *MockOrderRepository) InsertWithRetry(ctx context.Context, order *models.CombinedData) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) InsertOrder(order *models.CombinedData) error {
	args := m.Called(order)
	return args.Error(0)
}

func TestProcessMessage_ValidMessage(t *testing.T) {
	data := createValidData()

	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)

	msg := &kafka.Message{
		Topic: "orders",
		Value: jsonData,
	}

	ctx := context.Background()

	mockRepo := new(MockOrderRepository)
	mockRepo.On("InsertWithRetry", ctx, mock.AnythingOfType("*models.CombinedData")).Return(nil)

	err = processMessage(ctx, mockRepo, msg)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProcessMessage_InvalidJSON(t *testing.T) {
	msg := &kafka.Message{
		Value: []byte(`{invalid json}`),
	}

	err := processMessage(context.Background(), nil, msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestProcessMessage_ValidationFailed(t *testing.T) {
	data := createValidData()
	wrongEmail := ""
	data.Delivery.Email = &wrongEmail

	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)

	msg := &kafka.Message{
		Value: jsonData,
	}

	err = processMessage(context.Background(), nil, msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no address")
}
