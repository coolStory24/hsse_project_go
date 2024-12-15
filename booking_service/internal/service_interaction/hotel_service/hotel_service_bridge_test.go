package hotel_service_test

import (
	"booking_service/internal/service_interaction/hotel_service"
	"booking_service/internal/service_interaction/hotel_service/gen"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"testing"
)

type MockHotelServiceClient struct {
	mock.Mock
}

func (m *MockHotelServiceClient) GetHotelPrice(ctx context.Context, in *gen.GetHotelPriceRequest, opts ...grpc.CallOption) (*gen.GetHotelPriceResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*gen.GetHotelPriceResponse), args.Error(1)
}

func TestNewHotelServiceBridge(t *testing.T) {
	bridge, err := hotel_service.NewHotelServiceBridge("address")

	assert.Nil(t, err)
	assert.NotNil(t, bridge)
}

func TestHotelServiceBridge_GetHotelPrice(t *testing.T) {
	mockClient := new(MockHotelServiceClient)
	hotelBridge := &hotel_service.HotelServiceBridge{GrpcClient: mockClient}

	hotelId := uuid.New()
	expectedPrice := 100

	mockClient.On("GetHotelPrice", mock.Anything, &gen.GetHotelPriceRequest{HotelId: hotelId.String()}).Return(&gen.GetHotelPriceResponse{
		Price: int32(expectedPrice),
	}, nil)

	price, err := hotelBridge.GetHotelPrice(hotelId)

	assert.NoError(t, err)
	assert.Equal(t, expectedPrice, price)
	mockClient.AssertExpectations(t)
}

func TestHotelServiceBridge_GetHotelPrice_Error(t *testing.T) {
	mockClient := new(MockHotelServiceClient)
	hotelBridge := &hotel_service.HotelServiceBridge{GrpcClient: mockClient}

	hotelId := uuid.New()

	mockClient.On("GetHotelPrice", mock.Anything, &gen.GetHotelPriceRequest{HotelId: hotelId.String()}).Return(&gen.GetHotelPriceResponse{
		Price: int32(0),
	}, errors.New("failed to get price"))

	price, err := hotelBridge.GetHotelPrice(hotelId)

	assert.Error(t, err)
	assert.Equal(t, 0, price)
	mockClient.AssertExpectations(t)
}

func TestHotelServiceBridge_GetHotelPrice_Timeout(t *testing.T) {
	mockClient := new(MockHotelServiceClient)
	hotelBridge := &hotel_service.HotelServiceBridge{GrpcClient: mockClient}

	hotelId := uuid.New()

	mockClient.On("GetHotelPrice", mock.Anything, &gen.GetHotelPriceRequest{HotelId: hotelId.String()}).Return(
		&gen.GetHotelPriceResponse{
			Price: int32(0),
		}, errors.New("context deadline exceeded"))

	price, err := hotelBridge.GetHotelPrice(hotelId)

	assert.Error(t, err)
	assert.Equal(t, 0, price)
	mockClient.AssertExpectations(t)
}
