package server

import (
	"booking_service/internal/config"
	"booking_service/internal/rest"
	"booking_service/internal/services"
	"github.com/gorilla/mux"
)

func SetupApiRouter(cfg *config.ServerConfig, bookingService services.IBookingService) *mux.Router {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix(cfg.Prefix).Subrouter()
	apiRouter.HandleFunc("/rent", rest.CreateRentHandler(bookingService)).Methods("POST")
	apiRouter.HandleFunc("/rent/{rent_id}", rest.UpdateRentHandler(bookingService)).Methods("PUT")
	apiRouter.HandleFunc("/rent", rest.GetRentsHandler(bookingService)).Methods("GET")
	apiRouter.HandleFunc("/rent/{rent_id}", rest.GetRentByIDHandler(bookingService)).Methods("GET")

	return apiRouter
}
