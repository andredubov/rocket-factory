package auth

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/iam/internal/model"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

func (a *authService) Whoami(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error) {
	session, err := a.cache.Get(ctx, sessionUUID)
	if err != nil {
		logger.Error(ctx, "authService.Whoami", zap.Error(err))
		return nil, err
	}

	return session, nil
}
