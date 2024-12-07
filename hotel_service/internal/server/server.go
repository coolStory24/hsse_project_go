package server

import (
	"github.com/gorilla/mux"
	"hotel_service/internal/config"
	"hotel_service/internal/server/endpoints"
	"hotel_service/internal/services"
)

func SetupApiRouter(cfg *config.ServerConfig, hotelService services.IHotelService) *mux.Router {
	router := mux.NewRouter()

	// Setup API routes
	apiRouter := router.PathPrefix(cfg.Prefix).Subrouter()
	apiRouter.HandleFunc("/hotel/", endpoints.CreateHotelHandler(hotelService)).Methods("POST")
	apiRouter.HandleFunc("/hotel/{hotel_id}", endpoints.UpdateHotelHandler(hotelService)).Methods("PUT")
	apiRouter.HandleFunc("/hotel/{hotel_id}", endpoints.GetHotelHandler(hotelService)).Methods("GET")
	apiRouter.HandleFunc("/hotel", endpoints.GetAllHotelsHandler(hotelService)).Methods("GET")
	apiRouter.HandleFunc("/hotel/{hotel_id}", endpoints.DeleteHotelHandler(hotelService)).Methods("DELETE")

	return apiRouter
}
