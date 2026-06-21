package handlers

import (
	"net/http"
	"restaurant-menu-api/internal/db"
	"restaurant-menu-api/internal/utils"
	"strconv"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	results := make(map[string]string)
	results["in_use"] = strconv.Itoa(db.GetDB().Stats().InUse)
	results["idle"] = strconv.Itoa(db.GetDB().Stats().Idle)
	results["open_connections"] = strconv.Itoa(db.GetDB().Stats().OpenConnections)
	w.WriteHeader(http.StatusOK)
	utils.CreateResponse(w, http.StatusOK, results)
}
