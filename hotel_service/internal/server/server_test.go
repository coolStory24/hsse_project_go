package server_test

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"hotel_service/internal/config"
	"hotel_service/internal/dtos/requests"
	"hotel_service/internal/dtos/responses"
	"hotel_service/internal/server"
	"hotel_service/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// region Hotel Service Mock
type MockHotelService struct {
	mock.Mock
}

func (m *MockHotelService) Create(req requests.CreateHotelRequest) (uuid.UUID, error) {
	args := m.Called(req)
	id := uuid.New()
	return id, args.Error(0)
}

func (m *MockHotelService) Update(hotelID uuid.UUID, req requests.UpdateHotelRequest) error {
	args := m.Called(hotelID, req)
	return args.Error(0)
}

func (m *MockHotelService) GetByID(hotelID uuid.UUID, includePastRents bool) (*responses.GetHotelResponse, error) {
	args := m.Called(hotelID, includePastRents)
	return args.Get(0).(*responses.GetHotelResponse), args.Error(1)
}

func (m *MockHotelService) GetAllHotels(adminID *uuid.UUID) (*responses.GetHotelsResponse, error) {
	args := m.Called(adminID)
	return args.Get(0).(*responses.GetHotelsResponse), args.Error(1)
}

func (m *MockHotelService) DeleteHotel(hotelID uuid.UUID) error {
	args := m.Called(hotelID)
	return args.Error(0)
}

func (m *MockHotelService) ExistsById(id uuid.UUID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

// endregion

// region Test Helpers
func setupTestRouter(hotelService services.IHotelService) *mux.Router {
	return server.SetupApiRouter(&config.ServerConfig{Prefix: "/api"}, hotelService)
}

// endregion

// region Test Cases

func TestCreateHotel_CommonCase_Ok(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	reqBody := requests.CreateHotelRequest{
		HotelName:  "Test Hotel",
		NightPrice: 10000,
	}

	id := uuid.New()
	mockService.On("Create", reqBody).Return(nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/hotel/", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var returnedId uuid.UUID
	err := json.NewDecoder(rec.Body).Decode(&returnedId)
	assert.NoError(t, err)
	assert.Equal(t, id, returnedId)

	mockService.AssertExpectations(t)
}

func TestCreateHotel_ServiceError_Error(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	reqBody := requests.CreateHotelRequest{
		HotelName:  "Test Hotel",
		NightPrice: -100000,
	}

	mockService.On("Create", reqBody).Return(nil, new(error)).Once()

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/hotel/", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	mockService.AssertExpectations(t)
}

func TestUpdateHotel_CommonCase_Ok(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	hotelID := uuid.New()
	reqBody := requests.UpdateHotelRequest{
		HotelName:  "Updated Hotel",
		NightPrice: 12000,
	}

	mockService.On("ExistsById", hotelID).Return(true, nil)
	mockService.On("Update", hotelID, reqBody).Return(nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/hotel/"+hotelID.String(), bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateHotel_HotelDoesNotExist_ErrorNotFound(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	hotelID := uuid.New()
	reqBody := requests.UpdateHotelRequest{
		HotelName:  "Updated Hotel",
		NightPrice: 12000,
	}

	mockService.On("ExistsById", hotelID).Return(false, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/hotel/"+hotelID.String(), bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetHotel_HotelExists_Ok(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	hotelID := uuid.New()
	response := &responses.GetHotelResponse{
		Id:         hotelID,
		HotelName:  "Test Hotel",
		NightPrice: 10000,
	}

	mockService.On("GetByID", hotelID, true).Return(response, nil)

	req := httptest.NewRequest("GET", "/api/hotel/"+hotelID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resBody responses.GetHotelResponse
	err := json.NewDecoder(rec.Body).Decode(&resBody)
	assert.NoError(t, err)
	assert.Equal(t, response, &resBody)

	mockService.AssertExpectations(t)
}

func TestGetHotel_HotelDoesNotExist_ErrorNotFound(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	hotelID := uuid.New()
	mockService.On("GetByID", hotelID, true).Return(nil, nil)

	req := httptest.NewRequest("GET", "/api/hotel/"+hotelID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockService.AssertExpectations(t)
}

func TestGetAllHotels_CommonCase_Ok(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	response := &responses.GetHotelsResponse{
		Hotels: []responses.GetHotelResponse{
			{Id: uuid.New(), HotelName: "Hotel A", NightPrice: 10000},
			{Id: uuid.New(), HotelName: "Hotel B", NightPrice: 20000},
		},
	}

	mockService.On("GetAllHotels", nil).Return(response, nil)

	req := httptest.NewRequest("GET", "/api/hotel/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resBody responses.GetHotelsResponse
	err := json.NewDecoder(rec.Body).Decode(&resBody)
	assert.NoError(t, err)
	assert.Equal(t, response, &resBody)

	mockService.AssertExpectations(t)
}

func TestGetAllHotelsByAdmin_CommonCase_Ok(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	adminID := uuid.New()
	response := &responses.GetHotelsResponse{
		Hotels: []responses.GetHotelResponse{
			{Id: uuid.New(), HotelName: "Hotel A", NightPrice: 10000},
		},
	}

	mockService.On("GetAllHotels", &adminID).Return(response, nil)

	req := httptest.NewRequest("GET", "/api/hotel?admin="+adminID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resBody responses.GetHotelsResponse
	err := json.NewDecoder(rec.Body).Decode(&resBody)
	assert.NoError(t, err)
	assert.Equal(t, response, &resBody)

	mockService.AssertExpectations(t)
}

func TestDeleteHotel_HotelExists_NoContent(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	hotelID := uuid.New()

	mockService.On("ExistsById", hotelID).Return(true, nil)
	mockService.On("DeleteHotel", hotelID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/hotel/"+hotelID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteHotel_HotelDoesNotExist_NoError(t *testing.T) {
	mockService := new(MockHotelService)
	router := setupTestRouter(mockService)

	hotelID := uuid.New()

	mockService.On("ExistsById", hotelID).Return(false, nil)
	mockService.On("DeleteHotel", hotelID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/hotel/"+hotelID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}

// endregion
