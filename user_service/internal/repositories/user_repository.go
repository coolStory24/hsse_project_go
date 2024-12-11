package repositories

import (
	"database/sql"
	"fmt"
	"user_service/internal/db"
	"user_service/internal/user"
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
		SELECT u.id, u.user_name, u.email, u.password_hash, u.role
		FROM user u
		WHERE u.email = $1`

	row := s.Db.Connection.QueryRow(query, email)

	var user user.UserModel
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}
