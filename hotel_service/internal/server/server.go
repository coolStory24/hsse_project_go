package server

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"hotel_service/internal/config"
	"hotel_service/internal/metrics"
	"hotel_service/internal/server/endpoints"
	"hotel_service/internal/services"
)

func SetupApiRouter(cfg *config.ServerConfig, hotelService services.IHotelService) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	apiRouter := router.PathPrefix(cfg.Prefix).Subrouter()
	apiRouter.Use(metrics.MetricsMiddleware)

	apiRouter.HandleFunc("/hotel", endpoints.CreateHotelHandler(hotelService)).Methods("POST")
	apiRouter.HandleFunc("/hotel/{hotel_id}", endpoints.UpdateHotelHandler(hotelService)).Methods("PUT")
	apiRouter.HandleFunc("/hotel/{hotel_id}", endpoints.GetHotelHandler(hotelService)).Methods("GET")
	apiRouter.HandleFunc("/hotel", endpoints.GetAllHotelsHandler(hotelService)).Methods("GET")
	apiRouter.HandleFunc("/hotel/{hotel_id}", endpoints.DeleteHotelHandler(hotelService)).Methods("DELETE")

	return router
}
