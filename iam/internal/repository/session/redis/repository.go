package redis

import (
	"fmt"

	"github.com/andredubov/rocket-factory/platform/pkg/cache"
)

const cacheKeyPrefix = "session:"

type sessionsRepository struct {
	cache cache.RedisClient
}

func NewSessionsRepository(redisClient cache.RedisClient) *sessionsRepository {
	return &sessionsRepository{
		cache: redisClient,
	}
}

func (r *sessionsRepository) getCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}
