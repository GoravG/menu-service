package db

import (
	"database/sql"
	"fmt"
	"restaurant-menu-api/internal/config"
	"restaurant-menu-api/internal/logger"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	once sync.Once
	db   *sql.DB
)

func InitializeDB() {
	once.Do(func() {
		logger.Info("Initializing database")
		cfg := config.GetConfig()
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			logger.Error("Error opening database: " + err.Error())
			panic(err)
		}
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(10 * time.Minute)

		logger.Info("Pinging database")
		if err := db.Ping(); err != nil {
			logger.Error("Error pinging database: " + err.Error())
			panic(err)
		}
		logger.Info("Database pinged")

		CreateTablesIfNotExists(db)
		logger.Info("Database initialized")
	})
}

func GetDB() *sql.DB {
	InitializeDB()
	return db
}
