package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"foodApp/handlers"
	"foodApp/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockOrderRepo struct{}
type mockOrderRequestRepo struct{}

func (m *mockOrderRepo) Insert(ctx context.Context, o models.Order) (models.Order, error) {
	o.ID = "uuid-123"
	return o, nil
}

func (m *mockOrderRequestRepo) Insert(ctx context.Context, or models.OrderRequest) (models.OrderRequest, error) {
	or.ID = "uuid-req"
	return or, nil
}

// Success case
func TestCreateOrder(t *testing.T) {
	productRepo := &mockProductRepo{}
	orderRepo := &mockOrderRepo{}
	orderRequestRepo := &mockOrderRequestRepo{}
	coupons := map[string]int{"HAPPYHRS": 2}

	h := handlers.NewOrderHandler(productRepo, orderRepo, orderRequestRepo, coupons)

	body := []byte(`{
        "couponCode": "HAPPYHRS",
        "items": [{"productId": "10", "quantity": 1}]
    }`)
	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/order", h.Create).Methods("POST")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "uuid-123", resp["id"])
	assert.NotNil(t, resp["items"])
	assert.NotNil(t, resp["products"])
	_, hasCoupon := resp["couponCode"]
	assert.False(t, hasCoupon, "response should not contain couponCode")
}

// Invalid request
func TestCreateOrder_InvalidJson(t *testing.T) {
	h := handlers.NewOrderHandler(&mockProductRepo{}, &mockOrderRepo{}, &mockOrderRequestRepo{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBuffer([]byte(`{ bad json`)))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/order", h.Create).Methods("POST")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// missing items
func TestCreateOrder_ValidationError(t *testing.T) {
	h := handlers.NewOrderHandler(&mockProductRepo{}, &mockOrderRepo{}, &mockOrderRequestRepo{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBuffer([]byte(`{"items":[]}`)))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/order", h.Create).Methods("POST")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}
