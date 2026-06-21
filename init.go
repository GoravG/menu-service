package main

import (
	"fmt"
	"restaurant-menu-api/internal/config"
	"restaurant-menu-api/internal/db"
	"restaurant-menu-api/internal/logger"
)

func init() {
	initializeLogger()
	initializeConfig()
	initializeDB()
}

func initializeLogger() {
	fmt.Println("Initializing logger")
	logger.InitializeLogger()
	logger.Info("Logger initialized")
}

func initializeConfig() {
	logger.Info("Initializing config")
	config.InitializeConfig()
	logger.Info("Config initialized")
}

func initializeDB() {
	logger.Info("Initializing database")
	db.InitializeDB()
	logger.Info("Database initialized")
}
