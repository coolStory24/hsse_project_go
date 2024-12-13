package rest

import (
	custom_errors "booking_service/internal/errors"
	"booking_service/internal/rest/dtos/requests"
	"booking_service/internal/services"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"log/slog"
	"strconv"
)

func CreateRentHandler(service services.IBookingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling the rent creation handler")
		var req requests.CreateRentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			slog.Error("Invalid request body" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		if req.CheckOutDate.Before(req.CheckInDate) {
			http.Error(w, "Check-out date cannot be before check-in date", http.StatusBadRequest)
			slog.Error("Check-out date cannot be before check-in date" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		rentID, err := service.CreateRent(req)
		if err != nil {
			if errors.As(err, new(*custom_errors.ServiceBadRequestError)) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				slog.Error(err.Error() + strconv.Itoa(http.StatusBadRequest))
			} else {
				http.Error(w, "Failed to create rent", http.StatusInternalServerError)
				slog.Error("Failed to create rent" + strconv.Itoa(http.StatusInternalServerError))
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rentID); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response" + strconv.Itoa(http.StatusInternalServerError))
			return
		}
		slog.Info("The rent was successfully created")
		slog.Info("Rent ID: " + rentID.String())
	}
}

func UpdateRentHandler(service services.IBookingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling the rent update handler")
		var req requests.UpdateRentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			slog.Error("Invalid request body" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		vars := mux.Vars(r)
		rentID, err := uuid.Parse(vars["rent_id"])
		if err != nil {
			http.Error(w, "Invalid rent ID", http.StatusBadRequest)
			slog.Error("Invalid rent ID" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		if err := service.UpdateRent(rentID, req); err != nil {
			http.Error(w, "Failed to update rent", http.StatusInternalServerError)
			slog.Error("Failed to update rent" + strconv.Itoa(http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusNoContent)
		slog.Info("The rent was successfully updated")
		slog.Info("Rent ID: " + rentID.String())
	}
}

func GetRentsHandler(service services.IBookingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling the rents getting handler")
		queryParams := r.URL.Query()
		clientIDStr := queryParams.Get("client")
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
		hotelID, errHotel := parseUUID(hotelIDStr)
		fromDate, errFrom := parseTime(from)
		toDate, errTo := parseTime(to)

		if errClient != nil || errHotel != nil || errFrom != nil || errTo != nil {
			http.Error(w, "Invalid data (failed to parse)", http.StatusBadRequest)
			slog.Error("Invalid data (failed to parse)" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		filter := requests.RentFilter{
			ClientID: clientID,
			HotelID:  hotelID,
			FromDate: fromDate,
			ToDate:   toDate,
		}

		rents, err := service.GetRents(filter)
		if err != nil {
			http.Error(w, "Failed to fetch rents", http.StatusInternalServerError)
			slog.Error("Failed to fetch rents" + strconv.Itoa(http.StatusInternalServerError))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rents); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response" + strconv.Itoa(http.StatusInternalServerError))
			return
		}
		slog.Info("The rents was successfully got")
		slog.Info("Client ID: " + clientID.String())
		slog.Info("Hotel ID: " + hotelID.String())
		slog.Info("FromDate: " + fromDate.String())
		slog.Info("ToDate: " + toDate.String())
	}
}

func GetRentByIDHandler(service services.IBookingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling the rent getting handler")
		vars := mux.Vars(r)
		rentID, err := uuid.Parse(vars["rent_id"])
		if err != nil {
			http.Error(w, "Invalid rent ID", http.StatusBadRequest)
			slog.Error("Invalid rent ID" + strconv.Itoa(http.StatusBadRequest))
			return
		}

		rent, err := service.GetRentByID(rentID)
		if err != nil {
			http.Error(w, "Failed to fetch rent", http.StatusInternalServerError)
			slog.Error("Failed to fetch rent" + strconv.Itoa(http.StatusInternalServerError))
			return
		}

		if rent == nil {
			http.Error(w, "Rent not found", http.StatusNotFound)
			slog.Error("Rent not found" + strconv.Itoa(http.StatusNotFound))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rent); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response" + strconv.Itoa(http.StatusInternalServerError))
			return
		}
		slog.Info("The rent was successfully got")
		slog.Info("Rent ID: " + rentID.String())
	}
}
