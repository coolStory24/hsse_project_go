package server_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"notification_service/internal/models"
	"notification_service/internal/server"
	"notification_service/internal/service/send_notification"
	"testing"
	"time"
)

type MockNotificationSender struct {
	mock.Mock
}

func (m *MockNotificationSender) Send(notification models.NotificationData) error {
	args := m.Called(notification)
	return args.Error(0)
}

func TestNotificationRequestHandler_handleNotificationRequest_Success(t *testing.T) {
	mockSender := new(MockNotificationSender)

	notificationData := models.NotificationData{
		UserContactData: &models.UserContactData{
			Phone: "123",
			Email: "test@gmail.com",
		},
		RentData: &models.RentData{
			ID:           uuid.New(),
			ClientID:     uuid.New(),
			HotelID:      uuid.New(),
			NightPrice:   10000,
			CheckInDate:  time.Now(),
			CheckOutDate: time.Now(),
		},
	}

	mockSender.On("Send", notificationData).Return(nil).Once()

	handler := server.NewNotificationRequestHandler([]send_notification.INotificationSender{mockSender})

	handler.HandleNotificationRequest(notificationData)

	mockSender.AssertExpectations(t)
}

func TestNotificationRequestHandler_handleNotificationRequest_Error(t *testing.T) {
	mockSender := new(MockNotificationSender)

	notificationData := models.NotificationData{
		UserContactData: &models.UserContactData{
			Phone: "123",
			Email: "test@gmail.com",
		},
		RentData: &models.RentData{
			ID:           uuid.New(),
			ClientID:     uuid.New(),
			HotelID:      uuid.New(),
			NightPrice:   10000,
			CheckInDate:  time.Now(),
			CheckOutDate: time.Now(),
		},
	}

	mockSender.On("Send", notificationData).Return(errors.New("failed to send")).Once()

	handler := server.NewNotificationRequestHandler([]send_notification.INotificationSender{mockSender})

	handler.HandleNotificationRequest(notificationData)

	mockSender.AssertExpectations(t)
}

func TestNotificationRequestHandler_handleNotificationRequest_MultipleSenders(t *testing.T) {
	mockSender1 := new(MockNotificationSender)
	mockSender2 := new(MockNotificationSender)

	notificationData := models.NotificationData{
		UserContactData: &models.UserContactData{
			Phone: "123",
			Email: "test@gmail.com",
		},
		RentData: &models.RentData{
			ID:           uuid.New(),
			ClientID:     uuid.New(),
			HotelID:      uuid.New(),
			NightPrice:   10000,
			CheckInDate:  time.Now(),
			CheckOutDate: time.Now(),
		},
	}

	mockSender1.On("Send", notificationData).Return(nil).Once()
	mockSender2.On("Send", notificationData).Return(errors.New("failed to send")).Once()

	handler := server.NewNotificationRequestHandler([]send_notification.INotificationSender{mockSender1, mockSender2})

	handler.HandleNotificationRequest(notificationData)

	mockSender1.AssertExpectations(t)
	mockSender2.AssertExpectations(t)
}
