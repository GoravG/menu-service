package services

import (
	"restaurant-menu-api/internal/logger"
	"restaurant-menu-api/internal/models"
)

func (s *Service) GetMenuCard() ([]models.MenuCardItem, error) {
	menuCardItems, err := s.menuCard.GetAll()
	if err != nil {
		logger.Error("Error getting menu card: " + err.Error())
		return nil, err
	}
	return menuCardItems, nil
}
