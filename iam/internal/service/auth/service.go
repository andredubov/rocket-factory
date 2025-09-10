package auth

import (
	"time"

	"github.com/andredubov/rocket-factory/iam/internal/service"
)

type authService struct {
	usersRepository service.UsersRepository
	cache           service.SessionsRepository
	cacheTTL        time.Duration
	passwordHasher  service.PasswordHasher
}

func NewAuthService(
	usersRepository service.UsersRepository,
	cache service.SessionsRepository,
	cacheTTL time.Duration,
	passwordHasher service.PasswordHasher,
) *authService {
	return &authService{
		usersRepository: usersRepository,
		cache:           cache,
		cacheTTL:        cacheTTL,
		passwordHasher:  passwordHasher,
	}
}
