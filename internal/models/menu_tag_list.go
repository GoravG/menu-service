package models

type MenuTagList struct {
	MenuItemName string `json:"menu_item_name"`
	Tag          Tag    `json:"tag"`
}

type MenuTagListRequest struct {
	MenuItemName string `json:"menu_item_name"`
	Tags         []Tag  `json:"tags"`
}
