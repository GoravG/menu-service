package app_test

import (
	"database/sql"
	"net/http"
	"testing"

	"restaurant-menu-api/internal/app"
	"restaurant-menu-api/internal/models"
	"restaurant-menu-api/internal/testutil"
)

func setupTest(t *testing.T) (http.Handler, *sql.DB) {
	t.Helper()

	database := testutil.SetupTestDB(t)
	router, err := app.NewRouter(database)
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}
	return router, database
}

func TestHealthCheck(t *testing.T) {
	router, _ := setupTest(t)

	rec := testutil.DoRequest(t, router, http.MethodGet, "/", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	response := testutil.DecodeResponse(t, rec)
	data := testutil.DecodeData[map[string]string](t, response)

	for _, key := range []string{"in_use", "idle", "open_connections"} {
		if _, ok := data[key]; !ok {
			t.Fatalf("expected health response to include %q", key)
		}
	}
}

func TestCreateCategoryAndGetCategories(t *testing.T) {
	router, _ := setupTest(t)

	createRec := testutil.DoRequest(t, router, http.MethodPost, "/categories", models.CategoryRequest{
		Name:        "Mains",
		Description: "Main course dishes",
	})
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, createRec.Code, createRec.Body.String())
	}

	getRec := testutil.DoRequest(t, router, http.MethodGet, "/categories", nil)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRec.Code)
	}

	response := testutil.DecodeResponse(t, getRec)
	categories := testutil.DecodeData[[]models.Category](t, response)

	if len(categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(categories))
	}
	if categories[0].Name != "Mains" {
		t.Fatalf("expected category name %q, got %q", "Mains", categories[0].Name)
	}
}

func TestCreateCategoryDuplicate(t *testing.T) {
	router, _ := setupTest(t)

	category := models.CategoryRequest{
		Name:        "Mains",
		Description: "Main course dishes",
	}

	firstRec := testutil.DoRequest(t, router, http.MethodPost, "/categories", category)
	if firstRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, firstRec.Code)
	}

	secondRec := testutil.DoRequest(t, router, http.MethodPost, "/categories", category)
	if secondRec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, secondRec.Code)
	}

	response := testutil.DecodeResponse(t, secondRec)
	message := testutil.DecodeData[string](t, response)
	if message != "category already exists" {
		t.Fatalf("expected duplicate error message, got %q", message)
	}
}

func TestCreateMenuItemRequiresCategory(t *testing.T) {
	router, _ := setupTest(t)

	rec := testutil.DoRequest(t, router, http.MethodPost, "/menu", models.MenuItemRequest{
		Name:         "Paneer Tikka",
		Description:  "Grilled cottage cheese",
		IsVegetarian: true,
		Available:    true,
		Category:     "Mains",
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	response := testutil.DecodeResponse(t, rec)
	message := testutil.DecodeData[string](t, response)
	if message != "category does not exist" {
		t.Fatalf("expected category validation error, got %q", message)
	}
}

func TestCreateMenuItemAndGetMenu(t *testing.T) {
	router, _ := setupTest(t)

	categoryRec := testutil.DoRequest(t, router, http.MethodPost, "/categories", models.CategoryRequest{
		Name:        "Mains",
		Description: "Main course dishes",
	})
	if categoryRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, categoryRec.Code)
	}

	menuItem := models.MenuItemRequest{
		Name:         "Paneer Tikka",
		Description:  "Grilled cottage cheese",
		IsVegetarian: true,
		Available:    true,
		Category:     "Mains",
	}

	createRec := testutil.DoRequest(t, router, http.MethodPost, "/menu", menuItem)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, createRec.Code, createRec.Body.String())
	}

	getRec := testutil.DoRequest(t, router, http.MethodGet, "/menu", nil)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRec.Code)
	}

	response := testutil.DecodeResponse(t, getRec)
	menuItems := testutil.DecodeData[[]models.MenuItem](t, response)

	if len(menuItems) != 1 {
		t.Fatalf("expected 1 menu item, got %d", len(menuItems))
	}
	if menuItems[0].Name != menuItem.Name {
		t.Fatalf("expected menu item name %q, got %q", menuItem.Name, menuItems[0].Name)
	}
}

func TestGetMenuCardWithPricesAndTags(t *testing.T) {
	router, database := setupTest(t)

	categoryRec := testutil.DoRequest(t, router, http.MethodPost, "/categories", models.CategoryRequest{
		Name:        "Mains",
		Description: "Main course dishes",
	})
	if categoryRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, categoryRec.Code)
	}

	tagRec := testutil.DoRequest(t, router, http.MethodPost, "/tags", models.TagRequest{
		Name:        "Spicy",
		Description: "Hot and spicy",
	})
	if tagRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, tagRec.Code)
	}

	menuItem := models.MenuItemRequest{
		Name:         "Paneer Tikka",
		Description:  "Grilled cottage cheese",
		IsVegetarian: true,
		Available:    true,
		Category:     "Mains",
	}
	menuRec := testutil.DoRequest(t, router, http.MethodPost, "/menu", menuItem)
	if menuRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, menuRec.Code)
	}

	priceRec := testutil.DoRequest(t, router, http.MethodPost, "/menu-price", models.MenuPriceRequest{
		MenuItemName: "Paneer Tikka",
		Price:        12.50,
		Currency:     "USD",
		PortionSize:  "HALF",
	})
	if priceRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, priceRec.Code, priceRec.Body.String())
	}

	fullPriceRec := testutil.DoRequest(t, router, http.MethodPost, "/menu-price", models.MenuPriceRequest{
		MenuItemName: "Paneer Tikka",
		Price:        18.00,
		Currency:     "USD",
		PortionSize:  "FULL",
	})
	if fullPriceRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, fullPriceRec.Code)
	}

	testutil.LinkMenuTag(t, database, "Paneer Tikka", "Spicy")

	menuCardRec := testutil.DoRequest(t, router, http.MethodGet, "/menu-card", nil)
	if menuCardRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, menuCardRec.Code, menuCardRec.Body.String())
	}

	response := testutil.DecodeResponse(t, menuCardRec)
	menuCardItems := testutil.DecodeData[[]models.MenuCardItem](t, response)

	if len(menuCardItems) != 1 {
		t.Fatalf("expected 1 menu card item, got %d", len(menuCardItems))
	}

	item := menuCardItems[0]
	if item.MenuItemName != "Paneer Tikka" {
		t.Fatalf("expected menu item name %q, got %q", "Paneer Tikka", item.MenuItemName)
	}
	if len(item.Prices) != 2 {
		t.Fatalf("expected 2 prices, got %d", len(item.Prices))
	}
	if item.Tags == nil || len(*item.Tags) != 1 || (*item.Tags)[0] != "Spicy" {
		t.Fatalf("expected tags [Spicy], got %#v", item.Tags)
	}
}

func TestGetMenuCardExcludesUnavailableItems(t *testing.T) {
	router, _ := setupTest(t)

	testutil.DoRequest(t, router, http.MethodPost, "/categories", models.CategoryRequest{
		Name:        "Mains",
		Description: "Main course dishes",
	})

	testutil.DoRequest(t, router, http.MethodPost, "/menu", models.MenuItemRequest{
		Name:         "Available Item",
		Description:  "On the menu",
		IsVegetarian: true,
		Available:    true,
		Category:     "Mains",
	})

	testutil.DoRequest(t, router, http.MethodPost, "/menu", models.MenuItemRequest{
		Name:         "Unavailable Item",
		Description:  "Off the menu",
		IsVegetarian: true,
		Available:    false,
		Category:     "Mains",
	})

	menuCardRec := testutil.DoRequest(t, router, http.MethodGet, "/menu-card", nil)
	if menuCardRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, menuCardRec.Code)
	}

	response := testutil.DecodeResponse(t, menuCardRec)
	menuCardItems := testutil.DecodeData[[]models.MenuCardItem](t, response)

	if len(menuCardItems) != 1 {
		t.Fatalf("expected 1 available menu card item, got %d", len(menuCardItems))
	}
	if menuCardItems[0].MenuItemName != "Available Item" {
		t.Fatalf("expected only available item, got %q", menuCardItems[0].MenuItemName)
	}
}
