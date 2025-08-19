package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kost0/L0/internal/models"
	"github.com/Kost0/L0/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderCache struct {
	mock.Mock
}

func (m *MockOrderCache) Get(orderID string) (*models.CombinedData, bool) {
	args := m.Called(orderID)
	return args.Get(0).(*models.CombinedData), args.Bool(1)
}

func (m *MockOrderCache) Set(orderID string, data *models.CombinedData) {
	m.Called(orderID, data)
}

func (m *MockOrderCache) WarmUpCache(db *sql.DB, repo repository.OrderRepository, ctx context.Context) error {
	args := m.Called(ctx, db, repo, ctx)
	return args.Error(0)
}

type MockSQLOrderRepository struct {
	mock.Mock
}

func (m *MockSQLOrderRepository) SelectWithRetry(ctx context.Context, orderID string) (*models.CombinedData, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(*models.CombinedData), args.Error(1)
}

func (m *MockSQLOrderRepository) SelectOrder(orderID string) (*models.CombinedData, error) {
	args := m.Called(orderID)
	return args.Get(0).(*models.CombinedData), args.Error(1)
}

func (m *MockSQLOrderRepository) InsertWithRetry(ctx context.Context, order *models.CombinedData) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockSQLOrderRepository) InsertOrder(order *models.CombinedData) error {
	args := m.Called(order)
	return args.Error(0)
}

func setupRouter(handler *Handler, orderID string) *httptest.ResponseRecorder {
	r := chi.NewRouter()
	r.Get("/orders/{orderID}", handler.GetOrderByID)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID, nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	return rr
}

func TestHandler_GetOrderByID_CacheHit(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepo := new(MockSQLOrderRepository)
	handler := &Handler{Repo: mockRepo, Cache: mockCache}

	orderID := "order-1"
	expectedData := &models.CombinedData{Order: models.Order{OrderUID: orderID}}

	mockCache.On("Get", orderID).Return(expectedData, true)

	rr := setupRouter(handler, orderID)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var response models.CombinedData
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, &response)
}

func TestHandler_GetOrderByID_CacheMiss_DBSuccess(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepoMock := new(MockSQLOrderRepository)
	handler := &Handler{Repo: mockRepoMock, Cache: mockCache}

	orderID := "order-1"
	expectedData := &models.CombinedData{Order: models.Order{OrderUID: orderID}}

	mockCache.On("Get", orderID).Return(&models.CombinedData{}, false)
	mockRepoMock.On("SelectWithRetry", mock.Anything, orderID).Return(expectedData, nil)

	handler.Repo = mockRepoMock
	mockCache.On("Set", orderID, expectedData).Return()

	rr := setupRouter(handler, orderID)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var response models.CombinedData
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, &response)

	mockCache.AssertExpectations(t)
	mockRepoMock.AssertExpectations(t)
}

func TestHandler_GetOrderByID_CacheMiss_DBNotFoud(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepo := new(MockSQLOrderRepository)
	handler := &Handler{Repo: mockRepo, Cache: mockCache}

	orderID := "order-unknown"

	mockCache.On("Get", orderID).Return(&models.CombinedData{}, false)
	mockRepo.On("SelectWithRetry", mock.Anything, orderID).Return(&models.CombinedData{}, sql.ErrNoRows)

	rr := setupRouter(handler, orderID)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Empty(t, rr.Body.String())

	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestHandler_GetOrderByID_CacheMiss_DBError(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepo := new(MockSQLOrderRepository)
	handler := &Handler{Repo: mockRepo, Cache: mockCache}

	orderID := "order-1"

	mockCache.On("Get", orderID).Return(&models.CombinedData{}, false)
	mockRepo.On("SelectWithRetry", mock.Anything, orderID).Return(&models.CombinedData{}, errors.New("db error"))

	rr := setupRouter(handler, orderID)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Empty(t, rr.Body.String())

	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
