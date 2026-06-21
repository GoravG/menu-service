package main

import (
	"fmt"
	"net/http"
	"restaurant-menu-api/internal/app"
	"restaurant-menu-api/internal/config"
	"restaurant-menu-api/internal/db"
	"restaurant-menu-api/internal/logger"
)

func main() {
	router, err := app.NewRouter(db.GetDB())
	if err != nil {
		logger.Error("Error initializing router: " + err.Error())
		panic(err)
	}

	addr := config.GetConfig().Host + ":" + fmt.Sprintf("%d", config.GetConfig().Port)
	logger.Info("Server started on " + addr)
	http.ListenAndServe(addr, router)
}
