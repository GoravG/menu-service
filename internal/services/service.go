package services

import (
	"database/sql"
	"errors"
	"restaurant-menu-api/internal/repository"
)

type Service struct {
	categories *repository.CategoryRepository
	menus      *repository.MenuRepository
	tags       *repository.TagRepository
	menuPrices *repository.MenuPriceRepository
}

func New(db *sql.DB) (*Service, error) {
	categories, err := repository.NewCategoryRepository(db)
	if err != nil {
		return nil, err
	}

	menus, err := repository.NewMenuRepository(db)
	if err != nil {
		return nil, err
	}

	tags, err := repository.NewTagRepository(db)
	if err != nil {
		return nil, err
	}

	menuPrices, err := repository.NewMenuPriceRepository(db)
	if err != nil {
		return nil, err
	}

	return &Service{
		categories: categories,
		menus:      menus,
		tags:       tags,
		menuPrices: menuPrices,
	}, nil
}

func isDuplicate(err error) bool {
	return errors.Is(err, repository.ErrDuplicate)
}
