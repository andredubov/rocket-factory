package redis

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/iam/internal/model"
	"github.com/andredubov/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/iam/internal/repository/model"
)

func (s *sessionsRepository) Get(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error) {
	cacheKey := s.getCacheKey(sessionUUID.String())

	values, err := s.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrSessionNotFound
	}

	var session repoModel.Session
	err = redigo.ScanStruct(values, &session)
	if err != nil {
		return nil, err
	}

	return converter.SessionToModel(session)
}
