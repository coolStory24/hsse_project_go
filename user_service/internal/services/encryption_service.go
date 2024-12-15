package services

import (
	"fmt"
	"user_service/internal/user"

	"golang.org/x/crypto/bcrypt"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type IEncryptionService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(storedHash, password string) bool
	GenerateToken(username string) (string, error)
	ParseToken(tokenStr string) (*JWTClaims, error)
}

type EncryptionService struct {
	key *string
}

type JWTClaims struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     user.Role `json:"role"`
	jwt.RegisteredClaims
}

func NewEncryptionService(encryptionKey string) *EncryptionService {
	return &EncryptionService{key: &encryptionKey}
}

func (e *EncryptionService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (e *EncryptionService) VerifyPassword(storedHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))

	return err == nil
}

func (e *EncryptionService) GenerateToken(id uuid.UUID, username string, email string, role user.Role) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &JWTClaims{
		Id:       id,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "user_service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(*e.key))
}

func (e *EncryptionService) ParseToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(*e.key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
