package server

import (
	"booking_service/internal/config"
	"booking_service/internal/metrics"
	"booking_service/internal/rest"
	"booking_service/internal/services"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupApiRouter(cfg *config.ServerConfig, bookingService services.IBookingService) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	apiRouter := router.PathPrefix(cfg.Prefix).Subrouter()
	apiRouter.Use(metrics.MetricsMiddleware)

	apiRouter.HandleFunc("/rent", rest.CreateRentHandler(bookingService)).Methods("POST")
	apiRouter.HandleFunc("/rent/{rent_id}", rest.UpdateRentHandler(bookingService)).Methods("PUT")
	apiRouter.HandleFunc("/rent", rest.GetRentsHandler(bookingService)).Methods("GET")
	apiRouter.HandleFunc("/rent/{rent_id}", rest.GetRentByIDHandler(bookingService)).Methods("GET")

	return router
}
