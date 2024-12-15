package requests

import (
	"user_service/internal/user"
)

type CreateRequest struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     user.Role `json:"role"`
	Password string    `json:"password"`
}
