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

type UserContactData struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type IUserServiceBridge interface {
	GetUserContactData(userId uuid.UUID) (*UserContactData, error)
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

func (u *UserServiceBridge) GetUserContactData(userId uuid.UUID) (*UserContactData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &gen.GetUserContactDataRequest{UserId: userId.String()}
	slog.Info("Sending request to get contact data of user with id " + userId.String())
	response, err := u.GrpcClient.GetUserContactDate(ctx, request)
	if err != nil {
		return nil, err
	}

	userContactData := &UserContactData{Email: response.Email, Phone: response.Phone}
	return userContactData, nil
}
