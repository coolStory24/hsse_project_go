package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"hotel_service/internal/config"
	"hotel_service/internal/dtos/requests"
	"hotel_service/internal/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func SetupApiRouter(cfg *config.ServerConfig, hotelService services.IHotelService) *mux.Router {
	router := mux.NewRouter()

	// Setup API routes
	apiRouter := router.PathPrefix(cfg.Prefix).Subrouter()
	apiRouter.HandleFunc("/hotel/", createHotelHandler(hotelService)).Methods("POST")
	apiRouter.HandleFunc("/hotel/{hotel_id}", updateHotelHandler(hotelService)).Methods("PUT")
	apiRouter.HandleFunc("/hotel/{hotel_id}", getHotelHandler(hotelService)).Methods("GET")
	apiRouter.HandleFunc("/hotel/", getAllHotelsHandler(hotelService)).Methods("GET")
	apiRouter.HandleFunc("/hotel/{hotel_id}", deleteHotelHandler(hotelService)).Methods("DELETE")

	return apiRouter
}

func NewServer(cfg *config.ServerConfig, hotelService services.IHotelService) {
	router := SetupApiRouter(cfg, hotelService)

	// Server configuration
	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	fmt.Printf("Server is starting on localhost%s\n", cfg.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Could not listen on %s: %v\n", cfg.Port, err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server exited gracefully")
}

// region Endpoints Handlers

func createHotelHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.CreateHotelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var id uuid.UUID
		id, err := service.Create(req)
		if err != nil {
			http.Error(w, "Failed to create hotel", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(id); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func updateHotelHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.UpdateHotelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		hotelID, err := uuid.Parse(vars["hotel_id"])
		if err != nil {
			http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
			return
		}

		if exists, err := service.ExistsById(hotelID); err != nil || !exists {
			http.Error(w, "Hotel with given id does not exist", http.StatusNotFound)
			return
		}

		if err := service.Update(hotelID, req); err != nil {
			http.Error(w, "Failed to update hotel", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func getHotelHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hotelID, err := uuid.Parse(vars["hotel_id"])
		if err != nil {
			http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
			return
		}

		pastRents := r.URL.Query().Get("past_rents") == "true"

		res, err := service.GetByID(hotelID, pastRents)
		if err != nil {
			http.Error(w, "Failed to get hotel", http.StatusInternalServerError)
			return
		}

		if res == nil {
			http.Error(w, "Hotel does not exist", http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(res); err != nil {
			// Handle encoding error
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func getAllHotelsHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminID := r.URL.Query().Get("admin")

		var adminUUID *uuid.UUID
		if adminID != "" {
			id, err := uuid.Parse(adminID)
			if err != nil {
				http.Error(w, "Invalid admin ID", http.StatusBadRequest)
				return
			}
			adminUUID = &id
		}

		res, err := service.GetAllHotels(adminUUID)
		if err != nil {
			http.Error(w, "Failed to fetch hotels", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func deleteHotelHandler(service services.IHotelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hotelID, err := uuid.Parse(vars["hotel_id"])
		if err != nil {
			http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
			return
		}

		if err := service.DeleteHotel(hotelID); err != nil {
			http.Error(w, "Failed to delete hotel", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// endregion
