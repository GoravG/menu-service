package repository

import (
	"database/sql"
	"restaurant-menu-api/internal/models"
	"strings"
)

type MenuCardRepository struct {
	db         *sql.DB
	getAllStmt *sql.Stmt
}

func NewMenuCardRepository(db *sql.DB) (*MenuCardRepository, error) {
	query := `
	SELECT menu_items.name as menu_item_name,
	menu_items.description as description,
	menu_items.available as available,
	menu_items.is_vegetarian as is_vegetarian,
	menu_items.category as category,
	menu_price_lists.price as price,
	menu_price_lists.currency as currency,
	menu_price_lists.portion_size as portion_size,
	GROUP_CONCAT(DISTINCT tags.name ORDER BY tags.name) as tags
	FROM menu_items
	LEFT JOIN menu_price_lists ON menu_items.name = menu_price_lists.menu_item_name
	LEFT JOIN menu_tags_list ON menu_items.name = menu_tags_list.menu_item_name
	LEFT JOIN tags ON menu_tags_list.tag = tags.name
	WHERE menu_items.available = true
	GROUP BY menu_items.name, menu_price_lists.portion_size, menu_price_lists.price, menu_price_lists.currency
	ORDER BY menu_items.name, menu_price_lists.portion_size
	`
	getAllStmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return &MenuCardRepository{
		db:         db,
		getAllStmt: getAllStmt,
	}, nil
}

func (r *MenuCardRepository) GetAll() ([]models.MenuCardItem, error) {
	rows, err := r.getAllStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemIndex := make(map[string]int)
	menuCardItems := []models.MenuCardItem{}

	for rows.Next() {
		var name, description, category string
		var available, isVegetarian bool
		var price sql.NullFloat64
		var currency, portionSize, tagsConcat sql.NullString

		if err := rows.Scan(&name, &description, &available, &isVegetarian, &category, &price, &currency, &portionSize, &tagsConcat); err != nil {
			return nil, err
		}

		idx, exists := itemIndex[name]
		if !exists {
			item := models.MenuCardItem{
				MenuItemName: name,
				Description:  description,
				Available:    available,
				IsVegetarian: isVegetarian,
				Category:     category,
				Prices:       []models.MenuCardPrice{},
			}
			if tagsConcat.Valid && tagsConcat.String != "" {
				tags := strings.Split(tagsConcat.String, ",")
				item.Tags = &tags
			}
			itemIndex[name] = len(menuCardItems)
			menuCardItems = append(menuCardItems, item)
			idx = itemIndex[name]
		}

		if price.Valid && currency.Valid && portionSize.Valid {
			menuCardItems[idx].Prices = append(menuCardItems[idx].Prices, models.MenuCardPrice{
				Price:       price.Float64,
				Currency:    currency.String,
				PortionSize: portionSize.String,
			})
		}
	}
	return menuCardItems, rows.Err()
}
