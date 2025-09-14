package handlers

import (
	"encoding/json"
	_interface "foodApp/app/product/interface"
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
	product, err := p.ProductRepository.Get(r.Context(), productID)
	if err != nil {
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(product)
}
