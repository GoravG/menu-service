package app

import (
	"database/sql"
	"net/http"
	"restaurant-menu-api/internal/handlers"
	"restaurant-menu-api/internal/middleware"
	"restaurant-menu-api/internal/services"
)

func NewRouter(database *sql.DB) (http.Handler, error) {
	service, err := services.New(database)
	if err != nil {
		return nil, err
	}

	handler := handlers.New(service)

	router := http.NewServeMux()
	router.Handle("/", middleware.Recover(handlers.HealthCheck(database)))
	router.Handle("GET /menu", middleware.Recover(http.HandlerFunc(handler.GetMenu)))
	router.Handle("POST /menu", middleware.Recover(http.HandlerFunc(handler.PostMenu)))
	router.Handle("GET /categories", middleware.Recover(http.HandlerFunc(handler.GetCategories)))
	router.Handle("POST /categories", middleware.Recover(http.HandlerFunc(handler.PostCategory)))
	router.Handle("GET /tags", middleware.Recover(http.HandlerFunc(handler.GetAllTags)))
	router.Handle("POST /tags", middleware.Recover(http.HandlerFunc(handler.PostTag)))
	router.Handle("GET /menu-price", middleware.Recover(http.HandlerFunc(handler.GetMenuPrice)))
	router.Handle("POST /menu-price", middleware.Recover(http.HandlerFunc(handler.PostMenuPrice)))
	router.Handle("GET /menu-card", middleware.Recover(http.HandlerFunc(handler.GetMenuCard)))

	return router, nil
}
