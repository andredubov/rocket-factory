package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/iam/internal/converter"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *AuthImplementation) Login(ctx context.Context, request *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	sessionUUID, err := a.authService.Login(ctx, request.Login, request.Password)
	if err != nil {
		logger.Error(ctx, "Login", zap.Error(err))
		return nil, err
	}

	return converter.SessionUUIDToLoginResponse(sessionUUID), nil
}
