package middleware

import (
	"net/http"
	"restaurant-menu-api/internal/logger"
	"restaurant-menu-api/internal/utils"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Error: " + err.(error).Error())
				utils.CreateResponse(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
