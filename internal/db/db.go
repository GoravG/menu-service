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

func getMenuItemsTableQuery() string {
	return `
	CREATE TABLE IF NOT EXISTS menu_items (
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		is_vegetarian BOOLEAN NOT NULL,
		available BOOLEAN DEFAULT TRUE,
		category VARCHAR(255) NOT NULL,
		PRIMARY KEY (name),
		FOREIGN KEY (category) REFERENCES categories(name)
	)
	`
}

func getTagsTableQuery() string {
	return `
	CREATE TABLE IF NOT EXISTS tags (
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		PRIMARY KEY (name)
	)
	`
}

func getCategoriesTableQuery() string {
	return `
	CREATE TABLE IF NOT EXISTS categories (
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		PRIMARY KEY (name)
	)
	`
}

func getMenuPriceListsTableQuery() string {
	return `
	CREATE TABLE IF NOT EXISTS menu_price_lists (
		menu_item_name VARCHAR(255) NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		currency VARCHAR(3) NOT NULL,
		portion_size VARCHAR(255) NOT NULL,
		PRIMARY KEY (menu_item_name, portion_size),
		FOREIGN KEY (menu_item_name) REFERENCES menu_items(name)
	)
	`
}

func getMenuTagsListTableQuery() string {
	return `
	CREATE TABLE IF NOT EXISTS menu_tags_list (
		menu_item_name VARCHAR(255) NOT NULL,
		tag VARCHAR(255) NOT NULL,
		PRIMARY KEY (menu_item_name, tag),
		FOREIGN KEY (menu_item_name) REFERENCES menu_items(name)
	)
	`
}

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
	createCategoryTableIfNotExists()
	createTagsTableIfNotExists()
	createMenuItemsTableIfNotExists()
	createMenuPriceListsTableIfNotExists()
	createMenuTagsListTableIfNotExists()
	logger.Info("Tables created if not exists")
}

func createCategoryTableIfNotExists() {
	query := getCategoriesTableQuery()
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating category table: " + err.Error())
		panic(err)
	}
}

func createTagsTableIfNotExists() {
	query := getTagsTableQuery()
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating tags table: " + err.Error())
		panic(err)
	}
}

func createMenuItemsTableIfNotExists() {
	query := getMenuItemsTableQuery()
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating menu items table: " + err.Error())
		panic(err)
	}
}

func createMenuPriceListsTableIfNotExists() {
	query := getMenuPriceListsTableQuery()
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating menu price lists table: " + err.Error())
		panic(err)
	}
}

func createMenuTagsListTableIfNotExists() {
	query := getMenuTagsListTableQuery()
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating menu tags list table: " + err.Error())
		panic(err)
	}
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
