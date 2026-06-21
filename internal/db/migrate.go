package db

import (
	"database/sql"
	"restaurant-menu-api/internal/logger"
)

func CreateTablesIfNotExists(database *sql.DB) {
	logger.Info("Creating tables if not exists")
	createCategoryTableIfNotExists(database)
	createTagsTableIfNotExists(database)
	createMenuItemsTableIfNotExists(database)
	createMenuPriceListsTableIfNotExists(database)
	createMenuTagsListTableIfNotExists(database)
	logger.Info("Tables created if not exists")
}

func createCategoryTableIfNotExists(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		PRIMARY KEY (name)
	)
	`
	if _, err := database.Exec(query); err != nil {
		logger.Error("Error creating category table: " + err.Error())
		panic(err)
	}
}

func createTagsTableIfNotExists(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS tags (
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		PRIMARY KEY (name)
	)
	`
	if _, err := database.Exec(query); err != nil {
		logger.Error("Error creating tags table: " + err.Error())
		panic(err)
	}
}

func createMenuItemsTableIfNotExists(database *sql.DB) {
	query := `
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
	if _, err := database.Exec(query); err != nil {
		logger.Error("Error creating menu items table: " + err.Error())
		panic(err)
	}
}

func createMenuPriceListsTableIfNotExists(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS menu_price_lists (
		menu_item_name VARCHAR(255) NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		currency VARCHAR(3) NOT NULL,
		portion_size VARCHAR(255) NOT NULL,
		PRIMARY KEY (menu_item_name, portion_size),
		FOREIGN KEY (menu_item_name) REFERENCES menu_items(name)
	)
	`
	if _, err := database.Exec(query); err != nil {
		logger.Error("Error creating menu price lists table: " + err.Error())
		panic(err)
	}
}

func createMenuTagsListTableIfNotExists(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS menu_tags_list (
		menu_item_name VARCHAR(255) NOT NULL,
		tag VARCHAR(255) NOT NULL,
		PRIMARY KEY (menu_item_name, tag),
		FOREIGN KEY (menu_item_name) REFERENCES menu_items(name)
	)
	`
	if _, err := database.Exec(query); err != nil {
		logger.Error("Error creating menu tags list table: " + err.Error())
		panic(err)
	}
}
