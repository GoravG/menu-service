package models

type Category struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
