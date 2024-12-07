package rest_test

import (
	"booking_service/dtos/requests"
	"booking_service/dtos/responses"
	"booking_service/internal/config"
	errors2 "booking_service/internal/errors"
	"booking_service/internal/server"
	"booking_service/internal/services"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// region Mock Booking Service
type MockBookingService struct {
	mock.Mock
}

func (m *MockBookingService) CreateRent(req requests.CreateRentRequest) (uuid.UUID, error) {
	args := m.Called(req)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockBookingService) UpdateRent(id uuid.UUID, req requests.UpdateRentRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockBookingService) GetRents(filter requests.RentFilter) (*responses.GetRentsResponse, error) {
	args := m.Called(filter)
	return args.Get(0).(*responses.GetRentsResponse), args.Error(1)
}

func (m *MockBookingService) GetRentByID(id uuid.UUID) (*responses.GetRentResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*responses.GetRentResponse), args.Error(1)
}

// endregion

// region Helpers
func setupTestRouter(service services.IBookingService) *mux.Router {
	return server.SetupApiRouter(&config.ServerConfig{Prefix: "/api"}, service)
}

// endregion

// region Tests

func TestCreateRent_CommonCase_Ok(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	checkInDate := time.Now().Truncate(time.Second)
	checkOutDate := checkInDate.Add(72 * time.Hour)
	reqBody := requests.CreateRentRequest{
		HotelID:      uuid.New(),
		ClientID:     uuid.New(),
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}

	rentID := uuid.New()
	mockService.On("CreateRent", reqBody).Return(rentID, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/rent", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var returnedID uuid.UUID
	err := json.NewDecoder(rec.Body).Decode(&returnedID)
	assert.NoError(t, err)
	assert.Equal(t, rentID, returnedID)
	mockService.AssertExpectations(t)
}

func TestCreateRent_InvalidBody_BadRequest(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	reqBody := []byte("{invalid_json}")

	req := httptest.NewRequest("POST", "/api/rent", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCreateRent_CheckOutLessThanCheckIn_ReturnBadRequest(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	checkInDate := time.Now().Truncate(time.Second)
	checkOutDate := checkInDate.Add(-72 * time.Hour)
	reqBody := requests.CreateRentRequest{
		HotelID:      uuid.New(),
		ClientID:     uuid.New(),
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/rent", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCreateRent_NotExistingHotel_ReturnNotFound(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	checkInDate := time.Now().Truncate(time.Second)
	checkOutDate := checkInDate.Add(72 * time.Hour)
	reqBody := requests.CreateRentRequest{
		HotelID:      uuid.New(),
		ClientID:     uuid.New(),
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}

	mockService.On("CreateRent", mock.Anything).Return(uuid.Nil, errors2.NewServiceBadRequestError("service error", ""))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/rent", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetRents_CommonCase_Ok(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	clientID := uuid.New()
	rentID := uuid.New()
	hotelierId := uuid.New()
	hotelId := uuid.New()
	checkInDate := time.Now().Truncate(time.Second)
	checkOutDate := checkInDate.Add(72 * time.Hour)

	rents := responses.GetRentsResponse{
		Rents: []responses.GetRentResponse{
			{
				ID:           rentID,
				HotelID:      hotelId,
				HotelierId:   hotelierId,
				ClientID:     clientID,
				NightPrice:   1000,
				CheckInDate:  checkInDate,
				CheckOutDate: checkOutDate,
			},
		},
	}

	mockService.On("GetRents", mock.Anything).Return(&rents, nil)

	request := fmt.Sprintf("/api/rent?client=%s&hotelier=%s&hotel=%s&from=%s&to=%s",
		clientID.String(), hotelierId.String(), hotelId.String(),
		url.QueryEscape(checkInDate.Format(time.RFC3339)),
		url.QueryEscape(checkInDate.Format(time.RFC3339)))

	req := httptest.NewRequest("GET", request, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resBody responses.GetRentsResponse
	err := json.NewDecoder(rec.Body).Decode(&resBody)
	assert.NoError(t, err)
	assert.Equal(t, rents, resBody)
	mockService.AssertExpectations(t)
}

func TestGetRents_CommonCaseWithSomeEmptyFields_Ok(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	clientID := uuid.New()

	rents := responses.GetRentsResponse{
		Rents: []responses.GetRentResponse{
			{
				ClientID: clientID,
			},
		},
	}

	mockService.On("GetRents", mock.Anything).Return(&rents, nil)

	request := fmt.Sprintf("/api/rent?client=%s", clientID.String())

	req := httptest.NewRequest("GET", request, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resBody responses.GetRentsResponse
	err := json.NewDecoder(rec.Body).Decode(&resBody)
	assert.NoError(t, err)
	assert.Equal(t, rents, resBody)
	mockService.AssertExpectations(t)
}

func TestGetRents_FilterOutAllRents_ReturnEmptyRentsArray(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	clientID := uuid.New()
	hotelierId := uuid.New()
	hotelId := uuid.New()
	checkInDate := time.Now().Truncate(time.Second)

	mockService.On("GetRents", mock.Anything).Return(&responses.GetRentsResponse{Rents: []responses.GetRentResponse{}}, nil)

	request := fmt.Sprintf("/api/rent?client=%s&hotelier=%s&hotel=%s&from=%s&to=%s",
		clientID.String(), hotelierId.String(), hotelId.String(),
		url.QueryEscape(checkInDate.Format(time.RFC3339)),
		url.QueryEscape(checkInDate.Format(time.RFC3339)))

	req := httptest.NewRequest("GET", request, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resBody responses.GetRentsResponse
	err := json.NewDecoder(rec.Body).Decode(&resBody)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(resBody.Rents))
	mockService.AssertExpectations(t)
}

func TestUpdateRentHandler_ValidRequest_Ok(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	checkInDate := time.Now().Truncate(time.Second)
	checkOutDate := checkInDate.Add(72 * time.Hour)

	rentID := uuid.New()
	updateRequest := requests.UpdateRentRequest{
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}
	reqBody, _ := json.Marshal(updateRequest)

	mockService.On("UpdateRent", rentID, updateRequest).Return(nil)

	req := httptest.NewRequest("PUT", "/api/rent/"+rentID.String(), bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateRentHandler_InvalidRequestBody_StatusBadRequest(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	reqBody := []byte("{invalid_json}")

	req := httptest.NewRequest("PUT", "/api/rent/some-rent-id", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateRentHandler_InvalidRentID_StatsBadRequest(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	// Valid request body
	checkInDate := time.Now().Truncate(time.Second)
	checkOutDate := checkInDate.Add(72 * time.Hour)

	updateRequest := requests.UpdateRentRequest{
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}
	reqBody, _ := json.Marshal(updateRequest)

	req := httptest.NewRequest("PUT", "/api/rent/invalid-rent-id", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetRentByIDHandler_ValidRentID(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	checkInDate := time.Now().Truncate(time.Second)
	rentID := uuid.New()
	expectedRent := &responses.GetRentResponse{
		ID:           rentID,
		HotelID:      uuid.New(),
		HotelierId:   uuid.New(),
		ClientID:     uuid.New(),
		NightPrice:   1000,
		CheckInDate:  checkInDate,
		CheckOutDate: checkInDate.Add(72 * time.Hour),
	}

	mockService.On("GetRentByID", rentID).Return(expectedRent, nil)

	req := httptest.NewRequest("GET", "/api/rent/"+rentID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resBody responses.GetRentResponse
	err := json.NewDecoder(rec.Body).Decode(&resBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedRent, &resBody)

	mockService.AssertExpectations(t)
}

func TestGetRentByIDHandler_InvalidRentID(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	req := httptest.NewRequest("GET", "/api/rent/invalid-id", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid rent ID")

	mockService.AssertExpectations(t)
}

func TestGetRentByIDHandler_RentNotFound(t *testing.T) {
	mockService := new(MockBookingService)
	router := setupTestRouter(mockService)

	rentID := uuid.New()
	mockService.On("GetRentByID", rentID).Return((*responses.GetRentResponse)(nil), nil)

	req := httptest.NewRequest("GET", "/api/rent/"+rentID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Rent not found")

	mockService.AssertExpectations(t)
}

// endregion
