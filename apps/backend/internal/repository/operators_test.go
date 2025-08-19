package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kost0/L0/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestSQLOrderRepository_InsertOrder_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SQLOrderRepository{DB: db}

	data := createValidData()

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO delivery").
		WithArgs(
			*data.Delivery.ID, *data.Delivery.Name, *data.Delivery.Phone,
			*data.Delivery.Zip, *data.Delivery.City, *data.Delivery.Address,
			*data.Delivery.Region, *data.Delivery.Email,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(
			data.Order.OrderUID, *data.Order.TrackNumber, *data.Order.Entry,
			*data.Order.DeliveryID, *data.Order.Locale,
			*data.Order.InternalSignature, *data.Order.CustomerID,
			*data.Order.DeliveryService, *data.Order.Shardkey, *data.Order.SmID,
			data.Order.DateCreated, *data.Order.OofShard,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO payment").
		WithArgs(
			*data.Payment.Transaction, *data.Payment.RequestID,
			*data.Payment.Currency, *data.Payment.Provider, *data.Payment.Amount,
			*data.Payment.PaymentDT, *data.Payment.Bank, *data.Payment.DeliveryCost,
			*data.Payment.GoodsTotal, *data.Payment.CustomFee,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	for _, item := range data.Items {
		mock.ExpectExec("INSERT INTO items").
			WithArgs(
				*item.ChrtID, *item.TrackNumber, *item.Price, *item.Rid,
				*item.Name, *item.Sale, *item.Size, *item.TotalPrice,
				*item.NmID, *item.Brand, *item.Status,
			).WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	err = repo.InsertOrder(data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLOrderRepository_InsertOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SQLOrderRepository{DB: db}

	data := createValidData()

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO delivery").WillReturnError(errors.New("db down"))

	mock.ExpectRollback()

	err = repo.InsertOrder(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db down")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLOrderRepository_InsertWithRetry_SuccessOnFirstAttempt(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SQLOrderRepository{DB: db}

	data := createValidData()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO delivery").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO payment").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO items").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.InsertOrder(data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLOrderRepository_InsertWithRetry_RetryOnConnDone(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SQLOrderRepository{DB: db}

	data := createValidData()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO delivery").WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO delivery").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO payment").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO items").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.InsertWithRetry(context.Background(), data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLOrderRepository_InsertWithRetry_NoRetryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SQLOrderRepository{DB: db}

	data := createValidData()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO delivery").WillReturnError(errors.New("invalid email"))
	mock.ExpectRollback()

	err = repo.InsertWithRetry(context.Background(), data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email")
}

func TestSQLOrderRepository_SelectOrder_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SQLOrderRepository{DB: db}

	orderID := "order-1"

	rowsOrder := sqlmock.NewRows([]string{"order_uid", "track_number", "entry", "delivery_id", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"}).
		AddRow(orderID, "WB", "WBIL", "del-1", "en", "", "cust", "meest", "9", 99, time.Now(), "1")

	mock.ExpectQuery("SELECT \\* FROM orders").WithArgs(orderID).WillReturnRows(rowsOrder)

	rowsDel := sqlmock.NewRows([]string{"id", "name", "phone", "zip", "city", "address", "region", "email"}).
		AddRow("del-1", "Test", "+7", "123", "City", "Addr", "Region", "test@com")

	mock.ExpectQuery("SELECT \\* FROM delivery").WithArgs("del-1").WillReturnRows(rowsDel)

	rowsPay := sqlmock.NewRows([]string{"transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"}).
		AddRow(orderID, "", "USD", "wb", 100, 123, "alpha", 50, 50, 0)

	mock.ExpectQuery("SELECT \\* FROM payment").WithArgs(orderID).WillReturnRows(rowsPay)

	rowsItems := sqlmock.NewRows([]string{"chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status"}).
		AddRow(123, "WB", 100, "rid", "Item", 0, "M", 100, 456, "Brand", 202)

	mock.ExpectQuery("SELECT \\* FROM items").WithArgs("WB").WillReturnRows(rowsItems)

	data, err := repo.SelectOrder(orderID)
	assert.NoError(t, err)
	assert.Equal(t, orderID, data.Order.OrderUID)
	assert.Equal(t, "test@com", *data.Delivery.Email)
	assert.Equal(t, orderID, *data.Payment.Transaction)
	assert.Len(t, data.Items, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}
