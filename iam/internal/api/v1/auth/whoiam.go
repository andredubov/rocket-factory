package api

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/iam/internal/converter"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *AuthImplementation) Whoami(ctx context.Context, request *auth_v1.WhoamiRequest) (*auth_v1.WhoamiResponse, error) {
	sessionUUID, err := uuid.Parse(request.SessionUuid)
	if err != nil {
		logger.Error(ctx, "AuthImplementation.Whoami", zap.Error(err))
		return nil, err
	}

	session, err := a.authService.Whoami(ctx, sessionUUID)
	if err != nil {
		logger.Error(ctx, "AuthImplementation.Whoami", zap.Error(err))
		return nil, err
	}

	return converter.SessionToWhoamiResponse(session)
}
