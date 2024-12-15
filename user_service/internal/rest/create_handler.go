package rest

import (
	"encoding/json"
	"net/http"
	"user_service/internal/dto/requests"
	"user_service/internal/services"
)

func NewCreateHandler(service services.IUserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.CreateRequest

		if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		createResponse, serviceErr := service.Create(req)

		if serviceErr != nil {
			http.Error(w, "Unable to create user", http.StatusBadRequest)
			return
		}

		jsonData, err := json.Marshal(createResponse)

		if err != nil {
			http.Error(w, "Unable to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(jsonData)
		w.Header().Set("Content-Type", "application/json")
	}
}
