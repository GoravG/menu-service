package models

type MenuPrice struct {
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	PortionSize  string  `json:"portion_size"`
	MenuItemName string  `json:"menu_item_name"`
}

type MenuPriceRequest struct {
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	PortionSize  string  `json:"portion_size"`
	MenuItemName string  `json:"menu_item_name"`
}
