package models

type MenuPrice struct {
	ID          int    `json:"id"`
	PriceListID int    `json:"price_list_id"`
	PortionSize string `json:"portion_size"`
	MenuItemID  int    `json:"menu_item_id"`
}

type MenuPriceRequest struct {
	PriceListID int    `json:"price_list_id"`
	PortionSize string `json:"portion_size"`
	MenuItemID  int    `json:"menu_item_id"`
}
