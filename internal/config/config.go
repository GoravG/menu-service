package config

import (
	"os"
	"restaurant-menu-api/internal/logger"
	"strconv"
	"sync"
)

type Config struct {
	DBPassword string
	DBUser     string
	DBName     string
	DBHost     string
	DBPort     string
	Port       int
	Host       string
}

var (
	config *Config
	once   sync.Once
)

func InitializeConfig() *Config {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Error("Error converting port to int: " + err.Error())
		panic(err)
	}
	once.Do(func() {
		config = &Config{
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBUser:     os.Getenv("DB_USER"),
			DBName:     os.Getenv("DB_NAME"),
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
			Port:       port,
			Host:       os.Getenv("HOST"),
		}
	})
	return config
}

func GetConfig() *Config {
	return InitializeConfig()
}
