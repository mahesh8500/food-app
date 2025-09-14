package handlers

import (
	"encoding/json"
	"errors"
	"foodApp/app/product/db"
	_interface "foodApp/app/product/interface"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type ProductHandler struct {
	ProductRepository _interface.Products
}

func NewProductHandler(repo _interface.Products) *ProductHandler {
	return &ProductHandler{ProductRepository: repo}
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	products, err := p.ProductRepository.GetList(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(products)
}

func (p *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productID := mux.Vars(r)["id"]

	if _, err := uuid.Parse(productID); err != nil {
		http.Error(w, "invalid ID supplied", http.StatusBadRequest)
		return
	}

	product, err := p.ProductRepository.Get(r.Context(), productID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(product)
}
