package handlers

import (
	"net/http"
	"restaurant-menu-api/internal/models"
	"restaurant-menu-api/internal/utils"
)

func (h *Handler) GetMenu(w http.ResponseWriter, r *http.Request) {
	menuItems := h.service.GetMenuItems()
	utils.CreateResponse(w, http.StatusOK, menuItems)
}

func (h *Handler) PostMenu(w http.ResponseWriter, r *http.Request) {
	var menuItem models.MenuItemRequest
	err := utils.ParseRequestBody(r, &menuItem)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	err = h.service.CreateMenuItem(menuItem)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CreateResponse(w, http.StatusCreated, "Menu item created successfully")
}

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories := h.service.GetAllCategories()
	utils.CreateResponse(w, http.StatusOK, categories)
}

func (h *Handler) PostCategory(w http.ResponseWriter, r *http.Request) {
	var category models.CategoryRequest
	err := utils.ParseRequestBody(r, &category)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	err = h.service.CreateCategory(category)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CreateResponse(w, http.StatusCreated, "Category created successfully")
}

func (h *Handler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags := h.service.GetAllTags()
	utils.CreateResponse(w, http.StatusOK, tags)
}

func (h *Handler) PostTag(w http.ResponseWriter, r *http.Request) {
	var tag models.TagRequest
	err := utils.ParseRequestBody(r, &tag)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	err = h.service.CreateTag(tag)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CreateResponse(w, http.StatusCreated, "Tag created successfully")
}

func (h *Handler) GetMenuPrice(w http.ResponseWriter, r *http.Request) {
	menuPrices := h.service.GetMenuPrices()
	utils.CreateResponse(w, http.StatusOK, menuPrices)
}

func (h *Handler) PostMenuPrice(w http.ResponseWriter, r *http.Request) {
	var menuPrice models.MenuPriceRequest
	err := utils.ParseRequestBody(r, &menuPrice)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	err = h.service.CreateMenuPrice(menuPrice)
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CreateResponse(w, http.StatusCreated, "Menu price created successfully")
}

func (h *Handler) GetMenuCard(w http.ResponseWriter, r *http.Request) {
	menuCardItems, err := h.service.GetMenuCard()
	if err != nil {
		utils.CreateResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.CreateResponse(w, http.StatusOK, menuCardItems)
}
