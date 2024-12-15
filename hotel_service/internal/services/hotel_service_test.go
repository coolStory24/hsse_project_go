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

func TestUpdateHotel_CommonCase_Ok(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	hotelService := services.NewHotelService(&Database{Connection: db})

	request1 := requests.CreateHotelRequest{
		HotelName:  "Test Hotel",
		NightPrice: 100,
		AdminId:    uuid.New(),
	}

	mock.ExpectExec("INSERT INTO hotels").
		WithArgs(sqlmock.AnyArg(), request1.HotelName, request1.NightPrice, request1.AdminId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var hotelID uuid.UUID

	hotelID, err := hotelService.Create(request1)

	assert.NoError(t, err)

	request2 := requests.UpdateHotelRequest{
		HotelName:  "Test Hotel",
		NightPrice: 150,
	}

	mock.ExpectExec("UPDATE hotels SET hotel_name = \\$1, night_price = \\$2 WHERE id = \\$3").
        WithArgs(request2.HotelName, request2.NightPrice, hotelID).
        WillReturnResult(sqlmock.NewResult(1, 1))

	err = hotelService.Update(hotelID, request2)

	assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteHotel_CommonCase_Ok(t *testing.T) {
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

	var hotelID uuid.UUID

	hotelID, err := hotelService.Create(request)

	assert.NoError(t, err)

	mock.ExpectExec("DELETE FROM hotels WHERE id = \\$1").
        WithArgs(hotelID).
        WillReturnResult(sqlmock.NewResult(1, 1))

	err = hotelService.DeleteHotel(hotelID)

	assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllHotels_CommonCase_Ok(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	hotelService := services.NewHotelService(&Database{Connection: db})
	adminID := uuid.New()
    rows := sqlmock.NewRows([]string{"id", "hotel_name", "night_price", "administrator_id"}).
        AddRow(uuid.New(), "Test Hotel 1", 100, adminID).
        AddRow(uuid.New(), "Test Hotel 2", 200, adminID)

    mock.ExpectQuery("SELECT id, hotel_name, night_price, administrator_id FROM hotels").
        WillReturnRows(rows)

    response, err := hotelService.GetAllHotels(&adminID)

    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Len(t, response.Hotels, 2)
    assert.Equal(t, "Test Hotel 1", response.Hotels[0].HotelName)
    assert.Equal(t, 100, response.Hotels[0].NightPrice)
    assert.Equal(t, adminID, response.Hotels[0].AdminId)
    assert.Equal(t, "Test Hotel 2", response.Hotels[1].HotelName)
    assert.Equal(t, 200, response.Hotels[1].NightPrice)
    assert.Equal(t, adminID, response.Hotels[1].AdminId)
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
