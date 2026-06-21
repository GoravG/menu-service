package handlers

import "restaurant-menu-api/internal/services"

type Handler struct {
	service *services.Service
}

func New(service *services.Service) *Handler {
	return &Handler{service: service}
}
