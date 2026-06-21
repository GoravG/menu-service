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
	err := utils.ParseRequestBody(r, &menuItem)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	services.CreateMenuItem(menuItem)
	utils.CreateResponse(w, http.StatusCreated, "Menu item created successfully")
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories := services.GetAllCategories()
	utils.CreateResponse(w, http.StatusOK, categories)
}

func PostCategory(w http.ResponseWriter, r *http.Request) {
	var category models.CategoryRequest
	err := utils.ParseRequestBody(r, &category)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	services.CreateCategory(category)
	utils.CreateResponse(w, http.StatusCreated, "Category created successfully")
}
