package services

import (
	"errors"
	"restaurant-menu-api/internal/logger"
	"restaurant-menu-api/internal/models"
)

func (s *Service) CreateCategory(category models.CategoryRequest) error {
	err := s.categories.Create(category)
	if err != nil {
		logger.Error("Error creating category: " + err.Error())
		if isDuplicate(err) {
			return errors.New("category already exists")
		}
		return errors.New("failed to create category")
	}
	return nil
}

func (s *Service) GetAllCategories() []models.Category {
	categories, err := s.categories.GetAll()
	if err != nil {
		logger.Error("Error getting categories: " + err.Error())
		return []models.Category{}
	}
	return categories
}
