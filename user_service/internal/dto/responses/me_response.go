package responses

import (
	"user_service/internal/user"

	"github.com/google/uuid"
)

type MeResponse struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     user.Role `json:"role"`
}
