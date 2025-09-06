package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/iam/internal/model"
)

type UsersRepository interface {
	Create(ctx context.Context, user model.User) error
	Get(ctx context.Context, filter model.Filter) (*model.User, error)
}

type SessionsRepository interface {
	Get(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error)
	Create(ctx context.Context, session model.Session) error
}

type PasswordHasher interface {
	HashAndSalt(plainPassword string) (string, error)
	ComparePasswords(hashedPassword, plainPassword string) error
}
