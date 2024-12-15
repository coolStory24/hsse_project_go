package send_notification_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"notification_service/internal/models"
	"notification_service/internal/service/send_notification"
	"testing"
)

type MockContentBuilder struct {
	mock.Mock
}

func (m *MockContentBuilder) BuildContent(notification models.NotificationData) string {
	args := m.Called(notification)
	return args.String(0)
}

func TestNewEmailSender(t *testing.T) {
	emailSender := send_notification.NewEmailSender(&MockContentBuilder{}, "smtp_server", 0, "username", "password")

	assert.NotNil(t, emailSender)
}
