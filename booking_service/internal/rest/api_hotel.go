package rest

import (
	"booking_service/dtos/requests"
	custom_errors "booking_service/internal/errors"
	"booking_service/internal/services"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func CreateRentHandler(service services.IBookingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.CreateRentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.CheckOutDate.Before(req.CheckInDate) {
			http.Error(w, "Check-out date cannot be before check-in date", http.StatusBadRequest)
			return
		}

		rentID, err := service.CreateRent(req)
		if err != nil {
			if errors.As(err, new(*custom_errors.ServiceBadRequestError)) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, "Failed to create rent", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rentID); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateRentHandler(service services.IBookingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.UpdateRentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		rentID, err := uuid.Parse(vars["rent_id"])
		if err != nil {
			http.Error(w, "Invalid rent ID", http.StatusBadRequest)
			return
		}

		if err := service.UpdateRent(rentID, req); err != nil {
			http.Error(w, "Failed to update rent", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetRentsHandler(service services.IBookingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		clientIDStr := queryParams.Get("client")
		hotelierIDStr := queryParams.Get("hotelier")
		hotelIDStr := queryParams.Get("hotel")
		from := queryParams.Get("from")
		to := queryParams.Get("to")

		// Helper function to parse UUID from string
		parseUUID := func(s string) (uuid.UUID, error) {
			if s == "" {
				return uuid.Nil, nil
			}
			id, err := uuid.Parse(s)
			if err != nil {
				return uuid.Nil, err
			}
			return id, nil
		}

		parseTime := func(s string) (*time.Time, error) {
			if s == "" {
				return nil, nil
			}
			timeValue, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return nil, err
			}
			return &timeValue, nil
		}

		clientID, errClient := parseUUID(clientIDStr)
		hotelierID, errHotelier := parseUUID(hotelierIDStr)
		hotelID, errHotel := parseUUID(hotelIDStr)
		fromDate, errFrom := parseTime(from)
		toDate, errTo := parseTime(to)

		if errClient != nil || errHotelier != nil || errHotel != nil || errFrom != nil || errTo != nil {
			http.Error(w, "Invalid data (failed to parse)", http.StatusBadRequest)
			return
		}

		filter := requests.RentFilter{
			ClientID:   clientID,
			HotelierID: hotelierID,
			HotelID:    hotelID,
			FromDate:   fromDate,
			ToDate:     toDate,
		}

		rents, err := service.GetRents(filter)
		if err != nil {
			http.Error(w, "Failed to fetch rents", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rents); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func GetRentByIDHandler(service services.IBookingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rentID, err := uuid.Parse(vars["rent_id"])
		if err != nil {
			http.Error(w, "Invalid rent ID", http.StatusBadRequest)
			return
		}

		rent, err := service.GetRentByID(rentID)
		if err != nil {
			http.Error(w, "Failed to fetch rent", http.StatusInternalServerError)
			return
		}

		if rent == nil {
			http.Error(w, "Rent not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rent); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
