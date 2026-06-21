package main

import (
	"fmt"
	"net/http"
	"restaurant-menu-api/internal/config"
	"restaurant-menu-api/internal/handlers"
	"restaurant-menu-api/internal/logger"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/health", handlers.HealthCheck)
	router.HandleFunc("GET /menu", handlers.GetMenu)
	router.HandleFunc("POST /menu", handlers.PostMenu)
	router.HandleFunc("GET /categories", handlers.GetCategories)
	router.HandleFunc("POST /categories", handlers.PostCategory)
	addr := config.GetConfig().Host + ":" + fmt.Sprintf("%d", config.GetConfig().Port)
	logger.Info("Server started on " + addr)
	http.ListenAndServe(addr, router)
}
