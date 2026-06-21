package services

import (
	"errors"
	"restaurant-menu-api/internal/logger"
	"restaurant-menu-api/internal/models"
)

func (s *Service) CreateMenuItem(menuItem models.MenuItemRequest) error {
	exists, err := s.categories.Exists(menuItem.Category)
	if err != nil {
		logger.Error("Error checking category: " + err.Error())
		return errors.New("failed to validate category")
	}
	if !exists {
		return errors.New("category does not exist")
	}

	err = s.menus.Create(menuItem)
	if err != nil {
		logger.Error("Error creating menu item: " + err.Error())
		if isDuplicate(err) {
			return errors.New("menu item already exists")
		}
		return errors.New("failed to create menu item")
	}
	return nil
}

func (s *Service) GetMenuItems() []models.MenuItem {
	menuItems, err := s.menus.GetAll()
	if err != nil {
		logger.Error("Error getting menu items: " + err.Error())
		return []models.MenuItem{}
	}
	return menuItems
}
