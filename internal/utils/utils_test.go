package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"restaurant-menu-api/internal/utils"
)

func TestParseRequestBodyValidJSON(t *testing.T) {
	body := strings.NewReader(`{"name":"Paneer Tikka","description":"Grilled cottage cheese","is_vegetarian":true,"available":true,"category":"Mains"}`)
	req := httptest.NewRequest("POST", "/menu", body)

	var payload struct {
		Name string `json:"name"`
	}
	err := utils.ParseRequestBody(req, &payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if payload.Name != "Paneer Tikka" {
		t.Fatalf("expected name %q, got %q", "Paneer Tikka", payload.Name)
	}
}

func TestParseRequestBodyInvalidJSON(t *testing.T) {
	body := strings.NewReader(`{invalid`)
	req := httptest.NewRequest("POST", "/menu", body)

	var payload struct {
		Name string `json:"name"`
	}
	err := utils.ParseRequestBody(req, &payload)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCreateResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	utils.CreateResponse(rec, http.StatusCreated, "created")

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var response map[string]json.RawMessage
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	var data string
	if err := json.Unmarshal(response["data"], &data); err != nil {
		t.Fatalf("failed to decode data field: %v", err)
	}
	if data != "created" {
		t.Fatalf("expected data %q, got %q", "created", data)
	}
}
