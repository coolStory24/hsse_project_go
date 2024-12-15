package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"user_service/internal/services"
)

func NewMeHandler(service services.IUserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		response, err := service.GetUserByToken(token)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(err.Error()))
			return
		}

		jsonData, err := json.Marshal(response)

		if err != nil {
			http.Error(w, "Unable to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(jsonData)
		w.Header().Set("Content-Type", "application/json")
	}
}
