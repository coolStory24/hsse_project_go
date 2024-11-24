package services

import (
	"hotel_service/internal/dtos/requests"
	"hotel_service/internal/dtos/responses"

	"github.com/google/uuid"
)

type IHotelService interface {
	Create(request requests.CreateHotelRequest) (uuid.UUID, error)
	Update(hotelID uuid.UUID, request requests.UpdateHotelRequest) error
	GetByID(hotelID uuid.UUID) (*responses.GetHotelResponse, error)
	ExistsById(id uuid.UUID) (bool, error)
	GetAllHotels(adminUUID *uuid.UUID) (*responses.GetHotelsResponse, error)
	DeleteHotel(hotelID uuid.UUID) error
}

type HotelService struct {
}

func (s *HotelService) Create(hotel requests.CreateHotelRequest) (uuid.UUID, error) {
	// todo
	return uuid.New(), nil
}

func (s *HotelService) Update(hotelID uuid.UUID, request requests.UpdateHotelRequest) error {
	// todo
	return nil
}

func (s *HotelService) GetByID(hotelID uuid.UUID) (*responses.GetHotelResponse, error) {
	// todo
	return &responses.GetHotelResponse{NightPrice: 100}, nil
}

func (s *HotelService) GetAllHotels(adminID *uuid.UUID) (*responses.GetHotelsResponse, error) {
	// todo
	return &responses.GetHotelsResponse{
		Hotels: []responses.GetHotelResponse{
			{Id: uuid.New(), HotelName: "Hotel A", NightPrice: 10000},
			{Id: uuid.New(), HotelName: "Hotel B", NightPrice: 20000}}}, nil
}

func (s *HotelService) DeleteHotel(hotelID uuid.UUID) error {
	// todo
	return nil
}

func (s *HotelService) ExistsById(hotelId uuid.UUID) (bool, error) {
	// todo
	return true, nil
}
