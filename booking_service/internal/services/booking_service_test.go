package services_test

import (
	"booking_service/dtos/requests"
	"booking_service/dtos/responses"
	db2 "booking_service/internal/db"
	"booking_service/internal/services"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	return db, mock
}

func TestCreateRent_CommonCase_Ok(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})
	rentID := uuid.New()
	request := requests.CreateRentRequest{
		HotelID:      uuid.New(),
		ClientID:     uuid.New(),
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(24 * time.Hour),
	}

	mock.ExpectQuery(`INSERT INTO bookings \(hotel_id, client_id, check_in_date, check_out_date\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
		WithArgs(request.HotelID, request.ClientID, request.CheckInDate, request.CheckOutDate).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(rentID))

	id, err := bookingService.CreateRent(request)

	assert.NoError(t, err)
	assert.Equal(t, rentID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateRent_ErrorCase_DBError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})

	request := requests.CreateRentRequest{
		HotelID:      uuid.New(),
		ClientID:     uuid.New(),
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(24 * time.Hour),
	}

	mock.ExpectQuery(`INSERT INTO bookings \(hotel_id, client_id, check_in_date, check_out_date\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
		WithArgs(request.HotelID, request.ClientID, request.CheckInDate, request.CheckOutDate).
		WillReturnError(fmt.Errorf("database error"))

	_, err := bookingService.CreateRent(request)

	assert.Error(t, err)
	assert.Equal(t, "failed to create rent: database error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRent_Success(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})

	rentID := uuid.New()
	request := requests.UpdateRentRequest{
		HotelID:      uuid.New(),
		ClientID:     uuid.New(),
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(24 * time.Hour),
	}

	mock.ExpectExec(`UPDATE bookings SET hotel_id = \$2, client_id = \$3, check_in_date = \$4, check_out_date = \$5 WHERE id = \$1`).
		WithArgs(rentID, request.HotelID, request.ClientID, request.CheckInDate, request.CheckOutDate).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate 1 row updated

	err := bookingService.UpdateRent(rentID, request)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRent_NoRowsAffected(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})

	rentID := uuid.New()
	request := requests.UpdateRentRequest{
		HotelID:      uuid.New(),
		ClientID:     uuid.New(),
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(24 * time.Hour),
	}

	mock.ExpectExec(`UPDATE bookings SET hotel_id = \$2, client_id = \$3, check_in_date = \$4, check_out_date = \$5 WHERE id = \$1`).
		WithArgs(rentID, request.HotelID, request.ClientID, request.CheckInDate, request.CheckOutDate).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := bookingService.UpdateRent(rentID, request)

	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("rent with ID %s not found", rentID), err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRent_DBError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})

	rentID := uuid.New()
	request := requests.UpdateRentRequest{
		HotelID:      uuid.New(),
		ClientID:     uuid.New(),
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(24 * time.Hour),
	}

	mock.ExpectExec(`UPDATE bookings SET hotel_id = \$2, client_id = \$3, check_in_date = \$4, check_out_date = \$5 WHERE id = \$1`).
		WithArgs(rentID, request.HotelID, request.ClientID, request.CheckInDate, request.CheckOutDate).
		WillReturnError(fmt.Errorf("database error"))

	err := bookingService.UpdateRent(rentID, request)

	assert.Error(t, err)
	assert.Equal(t, "failed to update rent: database error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRentByID_CommonCase_ReturnRent(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})
	rentID := uuid.New()
	expectedRent := responses.GetRentResponse{
		ID:           rentID,
		HotelID:      uuid.New(),
		ClientID:     uuid.New(),
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(48 * time.Hour),
	}

	mock.ExpectQuery("SELECT b.id, b.hotel_id, b.client_id, b.check_in_date, b.check_out_date FROM bookings b").
		WithArgs(rentID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "hotel_id", "client_id", "check_in_date", "check_out_date"}).
			AddRow(expectedRent.ID, expectedRent.HotelID, expectedRent.ClientID, expectedRent.CheckInDate, expectedRent.CheckOutDate))

	rent, err := bookingService.GetRentByID(rentID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRent, *rent)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRentByID_NotFound(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})
	rentID := uuid.New()

	mock.ExpectQuery(`SELECT b.id, b.hotel_id, b.client_id, b.check_in_date, b.check_out_date FROM bookings b WHERE b.id = \$1`).
		WithArgs(rentID).
		WillReturnError(sql.ErrNoRows)

	rent, err := bookingService.GetRentByID(rentID)

	assert.NoError(t, err)
	assert.Nil(t, rent)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRentByID_DBError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})
	rentID := uuid.New()

	mock.ExpectQuery(`SELECT b.id, b.hotel_id, b.client_id, b.check_in_date, b.check_out_date FROM bookings b WHERE b.id = \$1`).
		WithArgs(rentID).
		WillReturnError(fmt.Errorf("database error"))

	rent, err := bookingService.GetRentByID(rentID)

	assert.Error(t, err)
	assert.Nil(t, rent)
	assert.Equal(t, "failed to fetch rent: database error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRents_WithFullFilter_ReturnFilteredRents(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})

	clientID := uuid.New()
	hotelierID := uuid.New()
	hotelID := uuid.New()
	fromDate := time.Now().Add(-24 * time.Hour)
	toDate := time.Now().Add(24 * time.Hour)

	filter := requests.RentFilter{
		ClientID:   &clientID,
		HotelierID: &hotelierID,
		HotelID:    &hotelID,
		FromDate:   &fromDate,
		ToDate:     &toDate,
	}

	mockRentID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "hotel_id", "client_id", "check_in_date", "check_out_date"}).
		AddRow(mockRentID, hotelID, clientID, fromDate, toDate)

	mock.ExpectQuery(`SELECT b.id, b.hotel_id, b.client_id, b.check_in_date, b.check_out_date FROM bookings b WHERE 1=1 AND b.client_id = \$1 AND b.hotelier_id = \$2 AND b.hotel_id = \$3 AND b.check_in_date >= \$4 AND b.check_out_date <= \$5`).
		WithArgs(clientID, hotelierID, hotelID, fromDate, toDate).
		WillReturnRows(rows)

	rents, err := bookingService.GetRents(filter)

	assert.NoError(t, err)
	assert.Len(t, rents.Rents, 1)
	assert.Equal(t, mockRentID, rents.Rents[0].ID)
	assert.Equal(t, hotelID, rents.Rents[0].HotelID)
	assert.Equal(t, clientID, rents.Rents[0].ClientID)
	assert.Equal(t, fromDate, rents.Rents[0].CheckInDate)
	assert.Equal(t, toDate, rents.Rents[0].CheckOutDate)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRents_WithPartialFilter_ReturnRents(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	bookingService := services.NewBookingService(&db2.Database{Connection: db})
	filter := requests.RentFilter{
		ClientID: new(uuid.UUID),
	}
	mockRentID := uuid.New()

	mock.ExpectQuery("SELECT b.id, b.hotel_id, b.client_id, b.check_in_date, b.check_out_date FROM bookings b").
		WillReturnRows(sqlmock.NewRows([]string{"id", "hotel_id", "client_id", "check_in_date", "check_out_date"}).
			AddRow(mockRentID, uuid.New(), uuid.New(), time.Now(), time.Now().Add(24*time.Hour)))

	rents, err := bookingService.GetRents(filter)

	assert.NoError(t, err)
	assert.Len(t, rents.Rents, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}