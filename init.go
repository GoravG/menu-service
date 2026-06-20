package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"restaurant-menu-api/internal/config"
	"restaurant-menu-api/internal/db"
	"restaurant-menu-api/internal/logger"
)

func init() {
	loadEnv()
	initializeLogger()
	initializeConfig()
	initializeDB()
}

func loadEnv() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "warning: failed to load .env: %v\n", err)
	}
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
