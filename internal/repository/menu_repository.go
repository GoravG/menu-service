package repository

import (
	"database/sql"
	"restaurant-menu-api/internal/models"
)

type MenuRepository struct {
	db         *sql.DB
	insertStmt *sql.Stmt
	getAllStmt *sql.Stmt
}

func NewMenuRepository(db *sql.DB) (*MenuRepository, error) {
	insertStmt, err := db.Prepare(`
		INSERT INTO menu_items (name, description, is_vegetarian, available, category)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, err
	}

	getAllStmt, err := db.Prepare(`
		SELECT name, description, is_vegetarian, available, category
		FROM menu_items
	`)
	if err != nil {
		insertStmt.Close()
		return nil, err
	}

	return &MenuRepository{
		db:         db,
		insertStmt: insertStmt,
		getAllStmt: getAllStmt,
	}, nil
}

func (r *MenuRepository) Create(item models.MenuItemRequest) error {
	_, err := r.insertStmt.Exec(
		item.Name,
		item.Description,
		item.IsVegetarian,
		item.Available,
		item.Category,
	)
	return mapError(err)
}

func (r *MenuRepository) GetAll() ([]models.MenuItem, error) {
	rows, err := r.getAllStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	menuItems := []models.MenuItem{}
	for rows.Next() {
		menuItem := models.MenuItem{}
		if err := rows.Scan(
			&menuItem.Name,
			&menuItem.Description,
			&menuItem.IsVegetarian,
			&menuItem.Available,
			&menuItem.Category,
		); err != nil {
			return nil, err
		}
		menuItems = append(menuItems, menuItem)
	}
	return menuItems, rows.Err()
}
