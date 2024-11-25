package services_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	. "hotel_service/internal/db"
	"hotel_service/internal/dtos/requests"
	"hotel_service/internal/services"
	"testing"
)

func createMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	return db, mock
}

func TestCreateHotel_CommonCase_Ok(t *testing.T) {
	// Create a mock database connection
	db, mock := createMockDB(t)
	defer db.Close()

	hotelService := services.NewHotelService(&Database{Connection: db})

	request := requests.CreateHotelRequest{
		HotelName:  "Test Hotel",
		NightPrice: 100,
		AdminId:    uuid.New(),
	}

	mock.ExpectExec("INSERT INTO hotels").
		WithArgs(sqlmock.AnyArg(), request.HotelName, request.NightPrice, request.AdminId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := hotelService.Create(request)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotelById_CommonCase_ReturnHotel(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	hotelService := services.NewHotelService(&Database{Connection: db})

	hotelID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "hotel_name", "night_price", "administrator_id"}).
		AddRow(hotelID, "Test Hotel", 100, uuid.New())

	mock.ExpectQuery("SELECT id, hotel_name, night_price, administrator_id FROM hotels").
		WithArgs(hotelID).
		WillReturnRows(rows)

	response, err := hotelService.GetByID(hotelID)

	assert.NoError(t, err)
	assert.Equal(t, hotelID, response.Id)
	assert.Equal(t, "Test Hotel", response.HotelName)
	assert.Equal(t, 100, response.NightPrice)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotelById_HotelWithIdDoesNotExist_ReturnNil(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	hotelService := &services.HotelService{Db: &Database{Connection: db}}

	hotelID := uuid.New()

	mock.ExpectQuery("SELECT id, hotel_name, night_price, administrator_id FROM hotels").
		WithArgs(hotelID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	result, err := hotelService.GetByID(hotelID)

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExistsById_Exists_ReturnTrue(t *testing.T) {
	// Create a mock database connection
	db, mock := createMockDB(t)
	defer db.Close()

	hotelService := services.NewHotelService(&Database{Connection: db})

	hotelID := uuid.New()
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(hotelID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := hotelService.ExistsById(hotelID)

	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExistsById_DoesNotExist_ReturnFalse(t *testing.T) {
	// Create a mock database connection
	db, mock := createMockDB(t)
	defer db.Close()

	hotelService := services.NewHotelService(&Database{Connection: db})

	hotelID := uuid.New()
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(hotelID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := hotelService.ExistsById(hotelID)

	assert.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}
