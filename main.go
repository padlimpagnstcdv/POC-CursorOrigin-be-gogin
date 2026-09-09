package main

import (
	"log"

	"be-gin-go/internal/core/service"
	"be-gin-go/internal/data/repository"
	"be-gin-go/internal/presentation/http/handler"
	"be-gin-go/internal/presentation/http/router"
)

func main() {
	// Wiring dependency: data -> core -> presentation
	productRepo := repository.NewProductMemoryRepository()
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	r := router.New(productHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
