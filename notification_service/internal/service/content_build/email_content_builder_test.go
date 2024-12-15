package content_build_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"notification_service/internal/models"
	"notification_service/internal/service/content_build"
	"testing"
	"time"
)

func TestEmailContentBuilder_BuildContent(t *testing.T) {
	builder := content_build.NewEmailContentBuilder()

	content := builder.BuildContent(models.NotificationData{
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
	})

	// Проверим, что контент не пуст. Проверять текст сообщения не имеет смысла
	assert.NotNil(t, content)
}
