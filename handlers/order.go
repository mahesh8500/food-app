package handlers

import (
	"encoding/json"
	orderInterface "foodApp/app/order/interface"
	orderReqInterface "foodApp/app/order_requests/interface"
	productInterface "foodApp/app/product/interface"
	"foodApp/models"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type OrderHandler struct {
	ProductRepo      productInterface.Products
	OrderRepo        orderInterface.Orders
	OrderRequestRepo orderReqInterface.OrderRequests
	Validate         *validator.Validate
	Coupons          map[string]int
}

func NewOrderHandler(prodRepo productInterface.Products, orderRepo orderInterface.Orders, orderReqRepo orderReqInterface.OrderRequests, coupons map[string]int) *OrderHandler {
	return &OrderHandler{
		ProductRepo:      prodRepo,
		OrderRepo:        orderRepo,
		OrderRequestRepo: orderReqRepo,
		Validate:         validator.New(),
		Coupons:          coupons,
	}
}

// DTO for request
type OrderReq struct {
	CouponCode string
	Items      []OrderItemReq `json:"items" validate:"required,min=1,dive"`
}

type OrderItemReq struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		coupon   *string
		req      OrderReq
		items    []map[string]interface{}
		products []map[string]interface{}
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, "validation exception", http.StatusUnprocessableEntity)
		return
	}

	itemsJSON, _ := json.Marshal(req.Items)
	if req.CouponCode != "" {
		coupon = &req.CouponCode
	}
	rawReq := models.OrderRequest{
		Items:      itemsJSON,
		CouponCode: coupon,
	}

	if _, err := h.OrderRequestRepo.Insert(r.Context(), rawReq); err != nil {
		http.Error(w, "failed to save order request", http.StatusInternalServerError)
		return
	}

	if req.CouponCode != "" {
		if len(req.CouponCode) < 8 || len(req.CouponCode) > 10 {
			http.Error(w, "invalid coupon code length", http.StatusUnprocessableEntity)
			return
		}
		if count, ok := h.Coupons[req.CouponCode]; !ok || count < 2 {
			http.Error(w, "invalid coupon code", http.StatusUnprocessableEntity)
			return
		}
	}

	for _, item := range req.Items {
		p, err := h.ProductRepo.Get(r.Context(), item.ProductID)
		if err != nil {
			http.Error(w, "validation exception", http.StatusUnprocessableEntity)
			return
		}
		items = append(items, map[string]interface{}{
			"productId": item.ProductID,
			"quantity":  item.Quantity,
		})
		products = append(products, map[string]interface{}{
			"id":       p.ID,
			"name":     p.Name,
			"price":    p.Price,
			"category": p.Category,
		})
	}
	productsJSON, _ := json.Marshal(products)

	dbOrder := models.Order{
		Items:    itemsJSON,
		Products: productsJSON,
	}

	savedOrder, err := h.OrderRepo.Insert(r.Context(), dbOrder)
	if err != nil {
		http.Error(w, "failed to save order", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"id":       savedOrder.ID,
		"items":    items,
		"products": products,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
