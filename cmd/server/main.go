package main

import (
	"context"
	orderDB "foodApp/app/order/db"
	orderReqDB "foodApp/app/order_requests/db"
	productDB "foodApp/app/product/db"
	"foodApp/config"
	"foodApp/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	_ = godotenv.Load()

	pgConnStr := config.NewDBConfig().GetPgConnStr()

	pool, err := productDB.NewPool(context.Background(), pgConnStr)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	couponMap, err := handlers.LoadCoupons([]string{
		"couponbase1.txt",
		"couponbase2.txt",
		"couponbase3.txt",
	})

	productRepo := productDB.NewPGRepository(pool)
	orderRepo := orderDB.NewOrderRepository(pool)
	orderRequestRepo := orderReqDB.NewOrderRequestRepository(pool)

	productHandler := handlers.NewProductHandler(productRepo)
	orderHandler := handlers.NewOrderHandler(productRepo, orderRepo, orderRequestRepo, couponMap)

	if err != nil {
		log.Fatalf("failed to load coupons: %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/products", productHandler.GetProducts).Methods("GET")
	router.HandleFunc("/products/{id}", productHandler.GetProductByID).Methods("GET")
	router.Handle("/order", handlers.ApiKeyMiddleware("apitest", http.HandlerFunc(orderHandler.Create))).Methods("POST")

	log.Println("Server running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
