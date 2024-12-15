package user_service_test

import (
	"booking_service/internal/service_interaction/user_service"
	"booking_service/internal/service_interaction/user_service/gen"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"testing"
)

type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) GetUserContactDate(ctx context.Context, in *gen.GetUserContactDataRequest, opts ...grpc.CallOption) (*gen.GetUserContactDataResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*gen.GetUserContactDataResponse), args.Error(1)
}

func TestNewUserServiceBridge(t *testing.T) {
	bridge, err := user_service.NewUserServiceBridge("address")

	assert.Nil(t, err)
	assert.NotNil(t, bridge)
}

func TestUserServiceBridge_GetUserContactData(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	userBridge := &user_service.UserServiceBridge{GrpcClient: mockClient}

	userId := uuid.New()
	expectedEmail := "test@example.com"
	expectedPhone := "1234567890"

	mockClient.On("GetUserContactDate", mock.Anything, &gen.GetUserContactDataRequest{UserId: userId.String()}).Return(&gen.GetUserContactDataResponse{
		Email: expectedEmail,
		Phone: expectedPhone,
	}, nil)

	contactData, err := userBridge.GetUserContactData(userId)

	assert.NoError(t, err)
	assert.NotNil(t, contactData)
	assert.Equal(t, expectedEmail, contactData.Email)
	assert.Equal(t, expectedPhone, contactData.Phone)
	mockClient.AssertExpectations(t)
}

func TestUserServiceBridge_GetUserContactData_Error(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	userBridge := &user_service.UserServiceBridge{GrpcClient: mockClient}

	userId := uuid.New()

	mockClient.On("GetUserContactDate", mock.Anything, &gen.GetUserContactDataRequest{UserId: userId.String()}).Return(&gen.GetUserContactDataResponse{}, errors.New("failed to get user contact data"))

	contactData, err := userBridge.GetUserContactData(userId)

	assert.Error(t, err)
	assert.Nil(t, contactData)
	mockClient.AssertExpectations(t)
}

func TestUserServiceBridge_GetUserContactData_Timeout(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	userBridge := &user_service.UserServiceBridge{GrpcClient: mockClient}

	userId := uuid.New()

	mockClient.On("GetUserContactDate", mock.Anything, &gen.GetUserContactDataRequest{UserId: userId.String()}).Return(&gen.GetUserContactDataResponse{}, errors.New("context deadline exceeded"))

	contactData, err := userBridge.GetUserContactData(userId)

	assert.Error(t, err)
	assert.Nil(t, contactData)
	mockClient.AssertExpectations(t)
}
