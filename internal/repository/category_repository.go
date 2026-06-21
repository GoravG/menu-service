package repository

import (
	"database/sql"
	"errors"
	"restaurant-menu-api/internal/models"
)

type CategoryRepository struct {
	db         *sql.DB
	insertStmt *sql.Stmt
	existsStmt *sql.Stmt
	getAllStmt *sql.Stmt
}

func NewCategoryRepository(db *sql.DB) (*CategoryRepository, error) {
	insertStmt, err := db.Prepare(`INSERT INTO categories (name, description) VALUES (?, ?)`)
	if err != nil {
		return nil, err
	}

	existsStmt, err := db.Prepare(`SELECT 1 FROM categories WHERE name = ? LIMIT 1`)
	if err != nil {
		insertStmt.Close()
		return nil, err
	}

	getAllStmt, err := db.Prepare(`SELECT name, description FROM categories`)
	if err != nil {
		insertStmt.Close()
		existsStmt.Close()
		return nil, err
	}

	return &CategoryRepository{
		db:         db,
		insertStmt: insertStmt,
		existsStmt: existsStmt,
		getAllStmt: getAllStmt,
	}, nil
}

func (r *CategoryRepository) Create(category models.CategoryRequest) error {
	_, err := r.insertStmt.Exec(category.Name, category.Description)
	return mapError(err)
}

func (r *CategoryRepository) Exists(name string) (bool, error) {
	var one int
	err := r.existsStmt.QueryRow(name).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.getAllStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		category := models.Category{}
		if err := rows.Scan(&category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}
