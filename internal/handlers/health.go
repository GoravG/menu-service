package handlers

import (
	"database/sql"
	"net/http"
	"restaurant-menu-api/internal/utils"
	"strconv"
)

func HealthCheck(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		results := make(map[string]string)
		results["in_use"] = strconv.Itoa(database.Stats().InUse)
		results["idle"] = strconv.Itoa(database.Stats().Idle)
		results["open_connections"] = strconv.Itoa(database.Stats().OpenConnections)
		utils.CreateResponse(w, http.StatusOK, results)
	})
}
