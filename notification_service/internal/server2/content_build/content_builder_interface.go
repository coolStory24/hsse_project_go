package content_build

import "hotel_service/internal/models"

type IContentBuilder interface {
	BuildContent(booking models.BookingRequest) string
}
