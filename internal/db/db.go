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

		createTablesIfNotExists()
		logger.Info("Database initialized")
	})
}

func getDBInstance() *sql.DB {
	InitializeDB()
	return db
}

func ExecuteQuery(query string, args ...interface{}) (int64, error) {
	db := getDBInstance()
	rows, err := db.Exec(query, args...)
	if err != nil {
		logger.Error("Error executing query: " + err.Error())
		return 0, err
	}
	rowsInserted, err := rows.LastInsertId()
	if err != nil {
		logger.Error("Error getting last inserted id: " + err.Error())
		return 0, err
	}
	return rowsInserted, nil
}

func MustExecuteQuery(query string, args ...interface{}) (int64, error) {
	_, err := ExecuteQuery(query, args...)
	if err != nil {
		logger.Error("Error executing query: " + err.Error())
		panic(err)
	}
	return 0, nil
}

func MustBegin() *sql.Tx {
	db := getDBInstance()
	tx, err := db.Begin()
	if err != nil {
		logger.Error("Error beginning transaction: " + err.Error())
		panic(err)
	}
	return tx
}

func createTablesIfNotExists() {
	logger.Info("Creating tables if not exists")
	query := `
	CREATE TABLE IF NOT EXISTS menu_items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		category VARCHAR(255) NOT NULL,
		available BOOLEAN NOT NULL
	)
	`
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating tables: " + err.Error())
		panic(err)
	}
	logger.Info("Tables created if not exists")
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	db := getDBInstance()
	rows, err := db.Query(query, args...)
	if err != nil {
		logger.Error("Error querying: " + err.Error())
		return nil, err
	}
	return rows, nil
}
