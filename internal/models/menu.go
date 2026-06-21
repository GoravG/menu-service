package models

type MenuItem struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsVegetarian bool   `json:"is_vegetarian"`
	Available    bool   `json:"available"`
	Category     string `json:"category"`
}

type MenuItemRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsVegetarian bool   `json:"is_vegetarian"`
	Available    bool   `json:"available"`
	Category     string `json:"category"`
}
