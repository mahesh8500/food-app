package handlers_test

import (
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

// mock that satisfies _interface.Products
type mockProductRepo struct{}

func (m *mockProductRepo) GetList(ctx context.Context) ([]models.Product, error) {
	return []models.Product{
		{ID: "1", Category: "Waffle", Name: "Chicken Waffle", Price: 12.5},
	}, nil
}
func (m *mockProductRepo) Get(ctx context.Context, id string) (models.Product, error) {
	return models.Product{ID: id, Category: "Waffle", Name: "Chicken Waffle", Price: 12.5}, nil
}

func TestGetProducts(t *testing.T) {
	h := handlers.NewProductHandler(&mockProductRepo{})

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products", h.GetProducts).Methods("GET")
	router.ServeHTTP(rr, req)

	t.Logf("RAW RESPONSE: %s", rr.Body.String())

	assert.Equal(t, http.StatusOK, rr.Code)

	var products []models.Product
	err := json.Unmarshal(rr.Body.Bytes(), &products)
	assert.NoError(t, err)

	assert.Len(t, products, 1)
	assert.Equal(t, "Chicken Waffle", products[0].Name)
}

func TestGetProductByID(t *testing.T) {
	h := handlers.NewProductHandler(&mockProductRepo{})

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", h.GetProductByID).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var product models.Product
	err := json.Unmarshal(rr.Body.Bytes(), &product)
	assert.NoError(t, err)
	assert.Equal(t, "Chicken Waffle", product.Name)
}
