package services

import (
	"restaurant-menu-api/internal/db"
	"restaurant-menu-api/internal/logger"
	"restaurant-menu-api/internal/models"
)

func CreateMenuItem(menuItem models.MenuItemRequest) int64 {
	query := "INSERT INTO menu_items (name, description, is_vegetarian, available, category) VALUES (?, ?, ?, ?, ?)"
	args := []interface{}{menuItem.Name, menuItem.Description, menuItem.IsVegetarian, menuItem.Available, menuItem.Category}
	rowsInserted, err := db.MustExecuteQuery(query, args...)
	if err != nil {
		logger.Error("Error creating menu item: " + err.Error())
		panic(err)
	}
	return rowsInserted
}

func GetMenuItems() []models.MenuItem {
	query := "SELECT * FROM menu_items"
	rows, err := db.Query(query)
	if err != nil {
		logger.Error("Error getting menu items: " + err.Error())
		panic(err)
	}
	defer rows.Close()
	menuItems := []models.MenuItem{}
	for rows.Next() {
		menuItem := models.MenuItem{}
		err = rows.Scan(&menuItem.Name, &menuItem.Description, &menuItem.IsVegetarian, &menuItem.Available, &menuItem.Category)
		if err != nil {
			logger.Error("Error scanning rows: " + err.Error())
			panic(err)
		}
		menuItems = append(menuItems, menuItem)
	}
	return menuItems
}

func CreateCategory(category models.CategoryRequest) int64 {
	query := "INSERT INTO categories (name, description) VALUES (?, ?)"
	args := []interface{}{category.Name, category.Description}
	rowsInserted, err := db.MustExecuteQuery(query, args...)
	if err != nil {
		logger.Error("Error creating category: " + err.Error())
		panic(err)
	}
	return rowsInserted
}

func GetAllCategories() []models.Category {
	query := "SELECT * FROM categories"
	rows, err := db.Query(query)
	if err != nil {
		logger.Error("Error getting categories: " + err.Error())
		panic(err)
	}
	defer rows.Close()
	categories := []models.Category{}
	for rows.Next() {
		category := models.Category{}
		err = rows.Scan(&category.Name, &category.Description)
		if err != nil {
			logger.Error("Error scanning rows: " + err.Error())
			panic(err)
		}
		categories = append(categories, category)
	}
	return categories
}
