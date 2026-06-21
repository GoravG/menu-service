package models

type MenuCardItem struct {
	MenuItemName string          `json:"menu_item_name"`
	Description  string          `json:"description"`
	Available    bool            `json:"available"`
	IsVegetarian bool            `json:"is_vegetarian"`
	Category     string          `json:"category"`
	Prices       []MenuCardPrice `json:"prices"`
	Tags         *[]string       `json:"tags,omitempty"`
}

type MenuCardPrice struct {
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	PortionSize string  `json:"portion_size"`
}
