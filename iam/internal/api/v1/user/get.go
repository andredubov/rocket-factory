package api

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/iam/internal/converter"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	user_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/user/v1"
)

func (u *UserImplementation) GetUser(ctx context.Context, request *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {
	userUUID, err := uuid.Parse(request.GetUserUuid())
	if err != nil {
		logger.Error(ctx, "GetUser", zap.Error(err))
		return nil, err
	}

	user, err := u.userService.Get(ctx, userUUID)
	if err != nil {
		logger.Error(ctx, "GetUser", zap.Error(err))
		return nil, err
	}

	return converter.UserToGetUserResponse(user)
}
