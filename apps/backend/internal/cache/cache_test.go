package cache

import (
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kost0/L0/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) SelectOrder(orderID string) (*models.CombinedData, error) {
	args := m.Called(orderID)
	return args.Get(0).(*models.CombinedData), args.Error(1)
}

func TestOrderCache_SetAndGet(t *testing.T) {
	cache := NewOrderCache(10 * time.Second)

	id := "order-1"
	data := &models.CombinedData{
		Order: models.Order{OrderUID: id},
	}

	cache.Set("order-1", data)

	result, ok := cache.Get("order-1")

	assert.True(t, ok)
	assert.Equal(t, data, result)
}

func TestOrderCache_GetNotFound(t *testing.T) {
	cache := NewOrderCache(10 * time.Second)

	result, found := cache.Get("not-found")
	assert.False(t, found)
	assert.NotNil(t, result)
}

func TestOrderCache_TTLExpiry(t *testing.T) {
	cache := NewOrderCache(100 * time.Millisecond)

	id := "order-1"
	data := &models.CombinedData{
		Order: models.Order{OrderUID: id},
	}

	cache.Set("order-1", data)

	_, found := cache.Get("order-1")
	assert.True(t, found)

	time.Sleep(110 * time.Millisecond)

	_, found = cache.Get("order-1")
	assert.False(t, found)
}

func TestOrderCache_WarmUpCache_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"order_uid"}).AddRow("order-1").AddRow("order-2")
	mock.ExpectQuery(`SELECT order_uid FROM orders.*7 days`).WillReturnRows(rows)

	mockRepo := new(MockOrderRepository)

	id1 := "order-1"
	id2 := "order-2"

	combined1 := &models.CombinedData{Order: models.Order{OrderUID: id1}}
	combined2 := &models.CombinedData{Order: models.Order{OrderUID: id2}}

	mockRepo.On("SelectOrder", "order-1").Return(combined1, nil)
	mockRepo.On("SelectOrder", "order-2").Return(combined2, nil)

	log.SetOutput(ioutil.Discard)

	cache := NewOrderCache(10 * time.Second)

	err = cache.WarmUpCache(db, mockRepo, context.Background())
	assert.NoError(t, err)

	res1, found1 := cache.Get("order-1")
	res2, found2 := cache.Get("order-2")

	assert.True(t, found1)
	assert.True(t, found2)
	assert.Equal(t, combined1, res1)
	assert.Equal(t, combined2, res2)
	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestOrderCache_WarmUpCache_SelectOrderError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"order_uid"}).AddRow("order-1")
	mock.ExpectQuery(`SELECT order_uid FROM orders`).WillReturnRows(rows)

	mockRepo := new(MockOrderRepository)
	mockRepo.On("SelectOrder", "order-1").Return((*models.CombinedData)(nil), sql.ErrNoRows)

	cache := NewOrderCache(10 * time.Second)
	err = cache.WarmUpCache(db, mockRepo, context.Background())

	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestGetRecentOrders_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"order_uid"}).AddRow("order-1")
	mock.ExpectQuery(`SELECT order_uid FROM orders.*7 days`).WillReturnRows(rows)

	mockRepo := new(MockOrderRepository)

	id := "order-1"

	expected := &models.CombinedData{Order: models.Order{OrderUID: id}}
	mockRepo.On("SelectOrder", "order-1").Return(expected, nil)

	data, err := getRecentOrders(db, mockRepo, context.Background())

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, id, data[0].Order.OrderUID)
	assert.NoError(t, mock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
}

func TestGetRecentOrders_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT order_uid FROM orders`).WillReturnError(sql.ErrTxDone)

	mockRepo := new(MockOrderRepository)

	_, err = getRecentOrders(db, mockRepo, context.Background())
	assert.Equal(t, sql.ErrTxDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
