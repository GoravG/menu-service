package services

import (
	"errors"
	"restaurant-menu-api/internal/logger"
	"restaurant-menu-api/internal/models"
)

func (s *Service) CreateMenuPrice(menuPrice models.MenuPriceRequest) error {
	err := s.menuPrices.Create(menuPrice)
	if err != nil {
		logger.Error("Error creating menu price: " + err.Error())
		return errors.New("failed to create menu price")
	}
	return nil
}

func (s *Service) GetMenuPrices() []models.MenuPrice {
	menuPrices, err := s.menuPrices.GetAll()
	if err != nil {
		logger.Error("Error getting menu prices: " + err.Error())
		return []models.MenuPrice{}
	}
	return menuPrices
}
