package services

import (
	"booking_service/dtos/requests"
	"booking_service/dtos/responses"
	"booking_service/internal/db"
	"github.com/google/uuid"
)

type IBookingService interface {
	CreateRent(request requests.CreateRentRequest) (uuid.UUID, error)
	UpdateRent(rentID uuid.UUID, request requests.UpdateRentRequest) error
	GetRentByID(rentID uuid.UUID) (*responses.GetRentResponse, error)
	GetRents(filter requests.RentFilter) (*responses.GetRentsResponse, error)
}

type BookingService struct {
	Db *db.Database
}

func NewBookingService(database *db.Database) *BookingService {
	return &BookingService{Db: database}
}

func (s *BookingService) CreateRent(request requests.CreateRentRequest) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s *BookingService) UpdateRent(rentID uuid.UUID, request requests.UpdateRentRequest) error {
	return nil
}

func (s *BookingService) GetRentByID(rentID uuid.UUID) (*responses.GetRentResponse, error) {
	return nil, nil
}

func (s *BookingService) GetRents(filter requests.RentFilter) (*responses.GetRentsResponse, error) {
	return nil, nil
}
