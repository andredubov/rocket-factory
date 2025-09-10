package redis

import (
	"context"

	"github.com/andredubov/rocket-factory/iam/internal/model"
	"github.com/andredubov/rocket-factory/iam/internal/repository/converter"
)

func (s *sessionsRepository) Create(ctx context.Context, session model.Session) error {
	cacheKey := s.getCacheKey(session.UUID.String())
	repoSession, err := converter.SessionToRepoModel(session)
	if err != nil {
		return err
	}

	err = s.cache.HashSet(ctx, cacheKey, repoSession)
	if err != nil {
		return err
	}

	cacheTTL := session.ExpiresAt.Sub(session.CreatedAt)

	return s.cache.Expire(ctx, cacheKey, cacheTTL)
}
