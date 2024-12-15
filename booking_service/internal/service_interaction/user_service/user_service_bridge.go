package user_service

import (
	"booking_service/internal/service_interaction/user_service/gen"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type UserData struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Phone string    `json:"phone"`
}

type IUserServiceBridge interface {
	GetUserContactData(token string) (*UserData, error)
}

type UserServiceBridge struct {
	GrpcClient gen.UserServiceClient
}

func NewUserServiceBridge(grpcAddress string) (*UserServiceBridge, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := gen.NewUserServiceClient(conn)

	return &UserServiceBridge{GrpcClient: client}, nil
}

func (u *UserServiceBridge) GetUserContactData(token string) (*UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &gen.GetUserDataRequest{Token: token}
	slog.Info("Sending request to get contact data of user with id " + token)
	response, err := u.GrpcClient.GetUserContactData(ctx, request)
	if err != nil {
		return nil, err
	}

	userContactData := &UserData{Email: response.Email, Phone: response.Phone}
	return userContactData, nil
}
