package services

import (
	"booking_service/dtos/requests"
	"booking_service/dtos/responses"
	"booking_service/internal/db"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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
	var rentID uuid.UUID

	query := `
        INSERT INTO bookings (hotel_id, client_id, check_in_date, check_out_date)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err := s.Db.Connection.QueryRow(query, request.HotelID, request.ClientID, request.CheckInDate, request.CheckOutDate).Scan(&rentID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create rent: %w", err)
	}
	return rentID, nil
}

func (s *BookingService) UpdateRent(rentID uuid.UUID, request requests.UpdateRentRequest) error {
	query := `
		UPDATE bookings
		SET hotel_id = $2, client_id = $3, check_in_date = $4, check_out_date = $5
		WHERE id = $1`
	result, err := s.Db.Connection.Exec(query, rentID, request.HotelID, request.ClientID, request.CheckInDate, request.CheckOutDate)
	if err != nil {
		return fmt.Errorf("failed to update rent: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("rent with ID %s not found", rentID)
	}
	return nil
}

func (s *BookingService) GetRentByID(rentID uuid.UUID) (*responses.GetRentResponse, error) {
	query := `
		SELECT b.id, b.hotel_id, b.client_id, b.check_in_date, b.check_out_date
		FROM bookings b
		WHERE b.id = $1`
	row := s.Db.Connection.QueryRow(query, rentID)

	var rent responses.GetRentResponse
	if err := row.Scan(&rent.ID, &rent.HotelID, &rent.ClientID, &rent.CheckInDate, &rent.CheckOutDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch rent: %w", err)
	}

	//Также нужно отправить запрос, чтобы получить стоимость комнаты

	return &rent, nil
}

func (s *BookingService) GetRents(filter requests.RentFilter) (*responses.GetRentsResponse, error) {
	query := `
		SELECT b.id, b.hotel_id, b.client_id, b.check_in_date, b.check_out_date
		FROM bookings b
		WHERE 1=1`

	params := []interface{}{}
	counter := 1

	if filter.ClientID != nil {
		query += fmt.Sprintf(" AND b.client_id = $%d", counter)
		params = append(params, *filter.ClientID)
		counter++
	}

	if filter.HotelierID != nil {
		query += fmt.Sprintf(" AND b.hotelier_id = $%d", counter)
		params = append(params, *filter.HotelierID)
		counter++
	}

	if filter.HotelID != nil {
		query += fmt.Sprintf(" AND b.hotel_id = $%d", counter)
		params = append(params, *filter.HotelID)
		counter++
	}

	if filter.FromDate != nil {
		query += fmt.Sprintf(" AND b.check_in_date >= $%d", counter)
		params = append(params, filter.FromDate)
		counter++
	}

	if filter.ToDate != nil {
		query += fmt.Sprintf(" AND b.check_out_date <= $%d", counter)
		params = append(params, filter.ToDate)
		counter++
	}

	rows, err := s.Db.Connection.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rents: %w", err)
	}
	defer rows.Close()

	var rents []responses.GetRentResponse
	for rows.Next() {
		var rent responses.GetRentResponse
		if err := rows.Scan(&rent.ID, &rent.HotelID, &rent.ClientID, &rent.CheckInDate, &rent.CheckOutDate); err != nil {
			return nil, fmt.Errorf("failed to scan rent: %w", err)
		}
		rents = append(rents, rent)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rents: %w", err)
	}

	// Нужно подгрузить для каждого отеля стоимости комнат

	return &responses.GetRentsResponse{Rents: rents}, nil
}
