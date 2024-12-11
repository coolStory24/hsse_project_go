package services

import (
	"user_service/internal/dto/requests"
	"user_service/internal/dto/responses"
	"user_service/internal/repositories"
)

type IUserService interface {
	Auth(request requests.AuthRequest) (responses.AuthResponse, error)
}

type UserService struct {
	Repository *repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{Repository: repository}
}

func (s *UserService) Auth(AuthRequest requests.AuthRequest) (responses.AuthResponse, error) {
	response := responses.AuthResponse{
		BearerToken: "Bearer some-token",
	}

	return response, nil
}
