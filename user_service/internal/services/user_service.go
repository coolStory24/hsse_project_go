package services

import (
	"fmt"
	"user_service/internal/dto/requests"
	"user_service/internal/dto/responses"
	"user_service/internal/repositories"
)

type IUserService interface {
	Auth(request requests.AuthRequest) (responses.AuthResponse, error)
	Create(request requests.CreateRequest) (responses.CreateResponse, error)
	GetUserByToken(token string) (responses.MeResponse, error)
}

type UserService struct {
	repository        *repositories.UserRepository
	encryptionService *EncryptionService
}

func NewUserService(repository *repositories.UserRepository, encryptionService *EncryptionService) *UserService {
	return &UserService{repository, encryptionService}
}

func (s *UserService) Auth(authRequest requests.AuthRequest) (responses.AuthResponse, error) {
	result, getErr := s.repository.GetByEmail(authRequest.Email)

	if getErr != nil || result == nil {
		return responses.AuthResponse{}, fmt.Errorf("auth failed")
	}

	isPasswordMatching := s.encryptionService.VerifyPassword(result.PasswordHash, authRequest.Password)

	response := responses.AuthResponse{}

	if isPasswordMatching {
		token, tokenErr := s.encryptionService.GenerateToken(result.Id, result.Username, result.Email, result.Role)

		if tokenErr != nil {
			return response, fmt.Errorf("auth failed")
		}

		response.BearerToken = token

		return response, nil
	} else {
		return response, fmt.Errorf("auth failed")
	}
}

func (s *UserService) Create(createRequest requests.CreateRequest) (responses.CreateResponse, error) {
	passwordHash, err := s.encryptionService.HashPassword(createRequest.Password)

	if err != nil {
		return responses.CreateResponse{}, err
	}

	result, createErr := s.repository.Create(createRequest.Username, createRequest.Email, createRequest.Role, passwordHash)

	if createErr != nil {
		return responses.CreateResponse{}, createErr
	}

	return responses.CreateResponse{Id: *result}, nil
}

func (s *UserService) GetUserByToken(token string) (responses.MeResponse, error) {
	result, err := s.encryptionService.ParseToken(token)

	if err != nil {
		return responses.MeResponse{}, fmt.Errorf("invalid token")
	}

	return responses.MeResponse{Id: result.Id, Username: result.Username, Email: result.Email, Role: result.Role}, nil
}
