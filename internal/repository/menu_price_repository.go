package repository

import (
	"database/sql"
	"errors"
	"restaurant-menu-api/internal/models"
)

type MenuPriceRepository struct {
	db         *sql.DB
	insertStmt *sql.Stmt
	existsStmt *sql.Stmt
	getAllStmt *sql.Stmt
}

func NewMenuPriceRepository(db *sql.DB) (*MenuPriceRepository, error) {
	insertStmt, err := db.Prepare(`INSERT INTO menu_price_lists (menu_item_name, price, currency, portion_size) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	existsStmt, err := db.Prepare(`SELECT 1 FROM menu_price_lists WHERE menu_item_name = ? AND portion_size = ? LIMIT 1`)
	if err != nil {
		insertStmt.Close()
		return nil, err
	}

	getAllStmt, err := db.Prepare(`SELECT menu_item_name, price, currency, portion_size FROM menu_price_lists`)
	if err != nil {
		insertStmt.Close()
		existsStmt.Close()
		return nil, err
	}

	return &MenuPriceRepository{
		db:         db,
		insertStmt: insertStmt,
		existsStmt: existsStmt,
		getAllStmt: getAllStmt,
	}, nil
}

func (r *MenuPriceRepository) Create(menuPrice models.MenuPriceRequest) error {
	_, err := r.insertStmt.Exec(menuPrice.MenuItemName, menuPrice.Price, menuPrice.Currency, menuPrice.PortionSize)
	return mapError(err)
}

func (r *MenuPriceRepository) Exists(menuItemName string, portionSize models.PortionSize) (bool, error) {
	var one int
	err := r.existsStmt.QueryRow(menuItemName, portionSize).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *MenuPriceRepository) GetAll() ([]models.MenuPrice, error) {
	rows, err := r.getAllStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	menuPrices := []models.MenuPrice{}
	for rows.Next() {
		menuPrice := models.MenuPrice{}
		if err := rows.Scan(&menuPrice.MenuItemName, &menuPrice.Price, &menuPrice.Currency, &menuPrice.PortionSize); err != nil {
			return nil, err
		}
		menuPrices = append(menuPrices, menuPrice)
	}
	return menuPrices, rows.Err()
}
