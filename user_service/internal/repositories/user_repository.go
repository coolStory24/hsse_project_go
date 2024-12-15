package repositories

import (
	"database/sql"
	"fmt"
	"user_service/internal/db"
	"user_service/internal/user"

	"github.com/google/uuid"
)

type IUserRepository interface {
	GetByEmail(email string) (user.UserModel, error)
}

type UserRepository struct {
	Db *db.Database
}

func NewUserRepository(database *db.Database) *UserRepository {
	return &UserRepository{Db: database}
}

func (s *UserRepository) GetByEmail(email string) (*user.UserModel, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.role
		FROM users u
		WHERE u.email = $1`

	row := s.Db.Connection.QueryRow(query, email)

	var user user.UserModel
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (s *UserRepository) Create(username string, email string, role user.Role, passwordHash string) (*uuid.UUID, error) {
	query := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var userId uuid.UUID

	err := s.Db.Connection.QueryRow(query, username, email, passwordHash, role).
		Scan(&userId)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &userId, nil
}
