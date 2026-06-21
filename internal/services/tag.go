package services

import (
	"errors"
	"restaurant-menu-api/internal/logger"
	"restaurant-menu-api/internal/models"
)

func (s *Service) CreateTag(tag models.TagRequest) error {
	err := s.tags.Create(tag)
	if err != nil {
		logger.Error("Error creating tag: " + err.Error())
		if isDuplicate(err) {
			return errors.New("tag already exists")
		}
		return errors.New("failed to create tag")
	}
	return nil
}

func (s *Service) GetAllTags() []models.Tag {
	tags, err := s.tags.GetAll()
	if err != nil {
		logger.Error("Error getting tags: " + err.Error())
		return []models.Tag{}
	}
	return tags
}
