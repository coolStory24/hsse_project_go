package rest

import (
	"encoding/json"
	"net/http"
	"user_service/internal/dto/requests"
	"user_service/internal/services"
)

func NewAuthHandler(service services.IUserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.AuthRequest

		if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		authResponse, serviceErr := service.Auth(req)

		if serviceErr != nil {
			http.Error(w, serviceErr.Error(), http.StatusBadRequest)
			return
		}

		jsonData, err := json.Marshal(authResponse)

		if err != nil {
			http.Error(w, "Unable to generate token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write(jsonData)
		w.Header().Set("Content-Type", "application/json")
	}
}
