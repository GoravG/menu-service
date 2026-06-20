package services

import (
	"restaurant-menu-api/internal/db"
	"restaurant-menu-api/internal/logger"
	"restaurant-menu-api/internal/models"
)

func CreateMenuItem(menuItem models.MenuItemRequest) int64 {
	query := "INSERT INTO menu_items (name, description, category, available) VALUES (?, ?, ?, ?)"
	args := []interface{}{menuItem.Name, menuItem.Description, menuItem.Category, menuItem.Available}
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
		err = rows.Scan(&menuItem.ID, &menuItem.Name, &menuItem.Description, &menuItem.Category, &menuItem.Available)
		if err != nil {
			logger.Error("Error scanning rows: " + err.Error())
			panic(err)
		}
		menuItems = append(menuItems, menuItem)
	}
	return menuItems
}
