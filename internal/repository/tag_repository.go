package repository

import (
	"database/sql"
	"restaurant-menu-api/internal/models"
)

type TagRepository struct {
	db         *sql.DB
	insertStmt *sql.Stmt
	getAllStmt *sql.Stmt
}

func NewTagRepository(db *sql.DB) (*TagRepository, error) {
	insertStmt, err := db.Prepare(`INSERT INTO tags (name, description) VALUES (?, ?)`)
	if err != nil {
		return nil, err
	}

	getAllStmt, err := db.Prepare(`SELECT name, description FROM tags`)
	if err != nil {
		insertStmt.Close()
		return nil, err
	}

	return &TagRepository{
		db:         db,
		insertStmt: insertStmt,
		getAllStmt: getAllStmt,
	}, nil
}

func (r *TagRepository) Create(tag models.TagRequest) error {
	_, err := r.insertStmt.Exec(tag.Name, tag.Description)
	return mapError(err)
}

func (r *TagRepository) GetAll() ([]models.Tag, error) {
	rows, err := r.getAllStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		tag := models.Tag{}
		if err := rows.Scan(&tag.Name, &tag.Description); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}
