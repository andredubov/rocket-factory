package api

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/iam/internal/model"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (uuid.UUID, error)
	Whoami(ctx context.Context, session uuid.UUID) (*model.Session, error)
}

type AuthImplementation struct {
	auth_v1.UnimplementedAuthServiceServer
	authService AuthService
}

func NewAuthImplementation(service AuthService) *AuthImplementation {
	return &AuthImplementation{
		authService: service,
	}
}
