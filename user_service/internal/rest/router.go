package rest

import (
	"user_service/internal/config"
	"user_service/internal/services"

	"github.com/gorilla/mux"
)

func SetupApiRouter(cfg *config.ServerConfig, userService services.IUserService) *mux.Router {
	router := mux.NewRouter()

	// Setup API routes
	apiRouter := router.PathPrefix(cfg.Prefix).Subrouter()
	apiRouter.HandleFunc("/user/auth", NewAuthHandler(userService)).Methods("POST")

	return apiRouter
}
