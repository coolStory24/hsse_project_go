package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"hotel_service/internal/config"
	"hotel_service/internal/dtos/requests"
	"hotel_service/internal/services"
	"log/slog"
	"net/http"
	"strconv"
)

func SetupApiRouter(cfg *config.ServerConfig, hotelService services.IHotelService) *mux.Router {
	router := mux.NewRouter()

	// Setup API routes
	apiRouter := router.PathPrefix(cfg.Prefix).Subrouter()
	apiRouter.HandleFunc("/hotel/", createHotelHandler(hotelService)).Methods("POST")
	apiRouter.HandleFunc("/hotel/{hotel_id}", updateHotelHandler(hotelService)).Methods("PUT")
	apiRouter.HandleFunc("/hotel/{hotel_id}", getHotelHandler(hotelService)).Methods("GET")
	apiRouter.HandleFunc("/hotel", getAllHotelsHandler(hotelService)).Methods("GET")
	apiRouter.HandleFunc("/hotel/{hotel_id}", deleteHotelHandler(hotelService)).Methods("DELETE")

	return apiRouter
}

// region Endpoints Handlers

func createHotelHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling the hotel creation handler")
		var req requests.CreateHotelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			slog.Error("Invalid request body" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		var id uuid.UUID
		id, err := service.Create(req)
		if err != nil {
			http.Error(w, "Failed to create hotel", http.StatusBadRequest)
			slog.Error("Failed to create hotel" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(id); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response" + strconv.Itoa(http.StatusInternalServerError))
			return
		}
		slog.Info("The hotel was successfully created")
		slog.Info("Hotel ID: " + id.String())
	}
}

func updateHotelHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling the hotel update handler")
		var req requests.UpdateHotelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			slog.Error("Invalid request body" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		vars := mux.Vars(r)
		hotelID, err := uuid.Parse(vars["hotel_id"])
		if err != nil {
			http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
			slog.Error("Invalid hotel ID" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		slog.Info("Hotel ID: " + hotelID.String())
		if exists, err := service.ExistsById(hotelID); err != nil || !exists {
			http.Error(w, "Hotel with given id does not exist", http.StatusNotFound)
			slog.Error("Hotel with given id does not exist" + strconv.Itoa(http.StatusNotFound))
			return
		}

		if err := service.Update(hotelID, req); err != nil {
			http.Error(w, "Failed to update hotel", http.StatusInternalServerError)
			slog.Error("Failed to update hotel" + strconv.Itoa(http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusNoContent)
		slog.Info("The hotel was successfully updated")
		slog.Info("Hotel ID: " + hotelID.String())
	}
}

func getHotelHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling the hotel getting handler")
		vars := mux.Vars(r)
		hotelID, err := uuid.Parse(vars["hotel_id"])
		if err != nil {
			http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
			slog.Error("Invalid hotel ID" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		slog.Info("Hotel ID: " + hotelID.String())
		res, err := service.GetByID(hotelID)
		if err != nil {
			http.Error(w, "Failed to get hotel", http.StatusInternalServerError)
			slog.Error("Failed to get hotel" + strconv.Itoa(http.StatusInternalServerError))
			return
		}

		if res == nil {
			http.Error(w, "Hotel does not exist", http.StatusNotFound)
			slog.Error("Hotel does not exist" + strconv.Itoa(http.StatusNotFound))
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(res); err != nil {
			// Handle encoding error
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response" + strconv.Itoa(http.StatusInternalServerError))
			return
		}
		slog.Info("The hotel was successfully got")
		slog.Info("Hotel ID: " + hotelID.String())
	}
}

func getAllHotelsHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling the all hotels getting handler")
		adminID := r.URL.Query().Get("admin")

		var adminUUID *uuid.UUID
		if adminID != "" {
			id, err := uuid.Parse(adminID)
			if err != nil {
				http.Error(w, "Invalid admin ID", http.StatusBadRequest)
				slog.Error("Invalid hotel ID" + strconv.Itoa(http.StatusBadRequest))
				return
			}
			adminUUID = &id
		}
		slog.Info("Admin ID: " + adminID)

		res, err := service.GetAllHotels(adminUUID)
		if err != nil {
			http.Error(w, "Failed to fetch hotels", http.StatusInternalServerError)
			slog.Error("Failed to fetch hotels" + strconv.Itoa(http.StatusInternalServerError))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response" + strconv.Itoa(http.StatusInternalServerError))
			return
		}
		slog.Info("The all hotels was successfully got")
		slog.Info("Admin ID: " + adminID)
	}
}

func deleteHotelHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling the hotel deletion handler")
		vars := mux.Vars(r)
		hotelID, err := uuid.Parse(vars["hotel_id"])
		if err != nil {
			http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
			slog.Error("Invalid hotel ID" + strconv.Itoa(http.StatusBadRequest))
			return
		}
		slog.Info("Hotel ID: " + hotelID.String())

		if err := service.DeleteHotel(hotelID); err != nil {
			http.Error(w, "Failed to delete hotel", http.StatusInternalServerError)
			slog.Error("Failed to delete hotel" + strconv.Itoa(http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusNoContent)
		slog.Info("The hotel was successfully deleted")
		slog.Info("Hotel ID: " + hotelID.String())
	}
}

// endregion
