package user_service_test

import (
	"booking_service/internal/service_interaction/user_service"
	"booking_service/internal/service_interaction/user_service/gen"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"testing"
)

type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) GetUserContactData(ctx context.Context, in *gen.GetUserDataRequest, opts ...grpc.CallOption) (*gen.GetUserDataResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*gen.GetUserDataResponse), args.Error(1)
}

func TestNewUserServiceBridge(t *testing.T) {
	bridge, err := user_service.NewUserServiceBridge("address")

	assert.Nil(t, err)
	assert.NotNil(t, bridge)
}

func TestUserServiceBridge_GetUserContactData(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	userBridge := &user_service.UserServiceBridge{GrpcClient: mockClient}

	userToken := "token"
	expectedEmail := "test@example.com"
	expectedPhone := "1234567890"

	mockClient.On("GetUserContactData", mock.Anything, &gen.GetUserDataRequest{Token: userToken}).Return(&gen.GetUserDataResponse{
		Email: expectedEmail,
		Phone: expectedPhone,
	}, nil)

	contactData, err := userBridge.GetUserContactData(userToken)

	assert.NoError(t, err)
	assert.NotNil(t, contactData)
	assert.Equal(t, expectedEmail, contactData.Email)
	assert.Equal(t, expectedPhone, contactData.Phone)
	mockClient.AssertExpectations(t)
}

func TestUserServiceBridge_GetUserContactData_Error(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	userBridge := &user_service.UserServiceBridge{GrpcClient: mockClient}

	userToken := "token"

	mockClient.On("GetUserContactData", mock.Anything, &gen.GetUserDataRequest{Token: userToken}).Return(&gen.GetUserDataResponse{}, errors.New("failed to get user contact data"))

	contactData, err := userBridge.GetUserContactData(userToken)

	assert.Error(t, err)
	assert.Nil(t, contactData)
	mockClient.AssertExpectations(t)
}

func TestUserServiceBridge_GetUserContactData_Timeout(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	userBridge := &user_service.UserServiceBridge{GrpcClient: mockClient}

	userToken := "token"

	mockClient.On("GetUserContactData", mock.Anything, &gen.GetUserDataRequest{Token: userToken}).Return(&gen.GetUserDataResponse{}, errors.New("context deadline exceeded"))

	contactData, err := userBridge.GetUserContactData(userToken)

	assert.Error(t, err)
	assert.Nil(t, contactData)
	mockClient.AssertExpectations(t)
}
