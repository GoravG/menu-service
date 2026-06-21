package models

type MenuPriceList struct {
	MenuItemName string      `json:"menu_item_name"`
	Price        float64     `json:"price"`
	Currency     string      `json:"currency"`
	PortionSize  PortionSize `json:"portion_size"`
}

type MenuPriceListRequest struct {
	MenuItemName string      `json:"menu_item_name"`
	Price        float64     `json:"price"`
	Currency     string      `json:"currency"`
	PortionSize  PortionSize `json:"portion_size"`
}

type PortionSize string

const (
	PortionSizeHalf PortionSize = "HALF"
	PortionSizeFull PortionSize = "FULL"
)
