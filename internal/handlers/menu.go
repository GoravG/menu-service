package handlers

import (
	"net/http"
	"restaurant-menu-api/internal/models"
	"restaurant-menu-api/internal/services"
	"restaurant-menu-api/internal/utils"
)

func GetMenu(w http.ResponseWriter, r *http.Request) {
	menuItems := services.GetMenuItems()
	utils.CreateResponse(w, http.StatusOK, menuItems)
}

func PostMenu(w http.ResponseWriter, r *http.Request) {
	var menuItem models.MenuItemRequest
	services.CreateMenuItem(menuItem)
	utils.CreateResponse(w, http.StatusCreated, "Menu item created successfully")
}
