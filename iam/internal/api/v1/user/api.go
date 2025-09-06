package api

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/iam/internal/model"
	user_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/user/v1"
)

type UserService interface {
	Register(ctx context.Context, userInfo model.UserInfo, password string) (*model.User, error)
	Get(ctx context.Context, uuid uuid.UUID) (*model.User, error)
}

type UserImplementation struct {
	user_v1.UnimplementedUserServiceServer
	userService UserService
}

func NewUserImplementation(service UserService) *UserImplementation {
	return &UserImplementation{
		userService: service,
	}
}
