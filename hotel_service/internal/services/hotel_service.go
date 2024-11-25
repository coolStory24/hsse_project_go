package services

import (
	"database/sql"
	"errors"
	"hotel_service/internal/db"
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
	Db *db.Database
}

func NewHotelService(database *db.Database) *HotelService {
	return &HotelService{Db: database}
}

func (s *HotelService) Create(request requests.CreateHotelRequest) (uuid.UUID, error) {
	hotelID := uuid.New()
	query := `INSERT INTO hotels (id, hotel_name, night_price, administrator_id) VALUES ($1, $2, $3, $4)`
	_, err := s.Db.Connection.Exec(query, hotelID, request.HotelName, request.NightPrice, request.AdminId)
	if err != nil {
		return uuid.Nil, err
	}
	return hotelID, nil
}

func (s *HotelService) Update(hotelID uuid.UUID, request requests.UpdateHotelRequest) error {
	// todo @svyatsharik
	return nil
}

func (s *HotelService) GetByID(hotelID uuid.UUID) (*responses.GetHotelResponse, error) {
	query := `SELECT id, hotel_name, night_price, administrator_id FROM hotels WHERE id = $1`
	row := s.Db.Connection.QueryRow(query, hotelID)

	var response responses.GetHotelResponse
	if err := row.Scan(&response.Id, &response.HotelName, &response.NightPrice, &response.AdminId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Hotel does not exist
		}
		return nil, err // some error occurred, but it's not about hotel does not exist
	}
	return &response, nil
}

func (s *HotelService) GetAllHotels(adminID *uuid.UUID) (*responses.GetHotelsResponse, error) {
	// todo @svyatsharik
	return &responses.GetHotelsResponse{
		Hotels: []responses.GetHotelResponse{
			{Id: uuid.New(), HotelName: "Hotel A", NightPrice: 10000},
			{Id: uuid.New(), HotelName: "Hotel B", NightPrice: 20000}}}, nil
}

func (s *HotelService) DeleteHotel(hotelID uuid.UUID) error {
	// todo svyatsharik
	return nil
}

func (s *HotelService) ExistsById(id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM hotels WHERE id = $1)`
	var exists bool
	err := s.Db.Connection.QueryRow(query, id).Scan(&exists)
	return exists, err
}
