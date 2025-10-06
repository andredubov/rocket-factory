package api

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
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
	authv3.UnimplementedAuthorizationServer
	authService AuthService
}

func NewAuthImplementation(service AuthService) *AuthImplementation {
	return &AuthImplementation{
		authService: service,
	}
}
