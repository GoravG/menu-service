package models

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
