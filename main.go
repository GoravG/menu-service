package main

import (
	"fmt"
	"net/http"
	"restaurant-menu-api/internal/config"
	"restaurant-menu-api/internal/db"
	"restaurant-menu-api/internal/handlers"
	"restaurant-menu-api/internal/logger"
	"restaurant-menu-api/internal/middleware"
	"restaurant-menu-api/internal/services"
)

func main() {
	service, err := services.New(db.GetDB())
	if err != nil {
		logger.Error("Error initializing services: " + err.Error())
		panic(err)
	}

	handler := handlers.New(service)

	router := http.NewServeMux()
	router.Handle("/", middleware.Recover(http.HandlerFunc(handlers.HealthCheck)))
	router.Handle("GET /menu", middleware.Recover(http.HandlerFunc(handler.GetMenu)))
	router.Handle("POST /menu", middleware.Recover(http.HandlerFunc(handler.PostMenu)))
	router.Handle("GET /categories", middleware.Recover(http.HandlerFunc(handler.GetCategories)))
	router.Handle("POST /categories", middleware.Recover(http.HandlerFunc(handler.PostCategory)))
	router.Handle("GET /tags", middleware.Recover(http.HandlerFunc(handler.GetAllTags)))
	router.Handle("POST /tags", middleware.Recover(http.HandlerFunc(handler.PostTag)))
	router.Handle("GET /menu-price", middleware.Recover(http.HandlerFunc(handler.GetMenuPrice)))
	router.Handle("POST /menu-price", middleware.Recover(http.HandlerFunc(handler.PostMenuPrice)))
	router.Handle("GET /menu-card", middleware.Recover(http.HandlerFunc(handler.GetMenuCard)))

	addr := config.GetConfig().Host + ":" + fmt.Sprintf("%d", config.GetConfig().Port)
	logger.Info("Server started on " + addr)
	http.ListenAndServe(addr, router)
}
