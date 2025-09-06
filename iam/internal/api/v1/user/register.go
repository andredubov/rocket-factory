package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/iam/internal/converter"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	user_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/user/v1"
)

func (u *UserImplementation) Register(ctx context.Context, request *user_v1.RegisterRequest) (*user_v1.RegisterResponse, error) {
	user := converter.UserFromRequest(request)

	user, err := u.userService.Register(ctx, user.Info, user.Password)
	if err != nil {
		logger.Error(ctx, "UserImplementation.Register", zap.Error(err))
		return nil, err
	}

	return converter.UserToRegisterResponse(user)
}
