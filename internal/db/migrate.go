package db

import "restaurant-menu-api/internal/logger"

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
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		PRIMARY KEY (name)
	)
	`
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating category table: " + err.Error())
		panic(err)
	}
}

func createTagsTableIfNotExists() {
	query := `
	CREATE TABLE IF NOT EXISTS tags (
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		PRIMARY KEY (name)
	)
	`
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating tags table: " + err.Error())
		panic(err)
	}
}

func createMenuItemsTableIfNotExists() {
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
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating menu items table: " + err.Error())
		panic(err)
	}
}

func createMenuPriceListsTableIfNotExists() {
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
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating menu price lists table: " + err.Error())
		panic(err)
	}
}

func createMenuTagsListTableIfNotExists() {
	query := `
	CREATE TABLE IF NOT EXISTS menu_tags_list (
		menu_item_name VARCHAR(255) NOT NULL,
		tag VARCHAR(255) NOT NULL,
		PRIMARY KEY (menu_item_name, tag),
		FOREIGN KEY (menu_item_name) REFERENCES menu_items(name)
	)
	`
	if _, err := db.Exec(query); err != nil {
		logger.Error("Error creating menu tags list table: " + err.Error())
		panic(err)
	}
}
