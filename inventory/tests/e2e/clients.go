//go:build integration

package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
	user_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/user/v1"
)

// TestClients содержит клиенты для всех сервисов
type TestClients struct {
	AuthClient      auth_v1.AuthServiceClient
	UserClient      user_v1.UserServiceClient
	InventoryClient inventory_v1.InventoryServiceClient
	iamConn         *grpc.ClientConn
	inventoryConn   *grpc.ClientConn
}

// NewTestClients создает подключения ко всем сервисам
func NewTestClients(ctx context.Context, env *TestEnvironment) (*TestClients, error) {
	// Создаем соединение с IAM сервисом
	iamConn, err := grpc.NewClient(
		env.IAMAppContainer.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial IAM service: %w", err)
	}

	// Создаем соединение с Inventory сервисом
	inventoryConn, err := grpc.NewClient(
		env.InventoryAppContainer.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		if closeErr := iamConn.Close(); closeErr != nil {
			logger.Error(ctx, "Failed to close IAM connection", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("failed to dial Inventory service: %w", err)
	}

	return &TestClients{
		AuthClient:      auth_v1.NewAuthServiceClient(iamConn),
		UserClient:      user_v1.NewUserServiceClient(iamConn),
		InventoryClient: inventory_v1.NewInventoryServiceClient(inventoryConn),
		iamConn:         iamConn,
		inventoryConn:   inventoryConn,
	}, nil
}

// Close закрывает все gRPC соединения
func (c *TestClients) Close() {
	ctx := context.Background()

	if c.iamConn != nil {
		if err := c.iamConn.Close(); err != nil {
			logger.Error(ctx, "Failed to close IAM connection", zap.Error(err))
		} else {
			logger.Debug(ctx, "IAM connection closed successfully")
		}
	}

	if c.inventoryConn != nil {
		if err := c.inventoryConn.Close(); err != nil {
			logger.Error(ctx, "Failed to close Inventory connection", zap.Error(err))
		} else {
			logger.Debug(ctx, "Inventory connection closed successfully")
		}
	}
}

// CloseWithContext закрывает соединения с контекстом для таймаута
func (c *TestClients) CloseWithContext(ctx context.Context) error {
	var errors []error

	if c.iamConn != nil {
		if err := c.iamConn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("IAM connection close failed: %w", err))
		}
	}

	if c.inventoryConn != nil {
		if err := c.inventoryConn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("inventory connection close failed: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple connection close errors: %v", errors)
	}

	return nil
}

// HealthCheck проверяет доступность всех сервисов
func (c *TestClients) HealthCheck(ctx context.Context) error {
	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Проверяем IAM service
	if _, err := c.AuthClient.Login(healthCtx, &auth_v1.LoginRequest{
		Login:    "healthcheck",
		Password: "healthcheck",
	}); err != nil {
		if status.Code(err) == codes.Unavailable {
			return fmt.Errorf("IAM service unavailable: %w", err)
		}
	}

	// Проверяем Inventory service
	if _, err := c.InventoryClient.GetPart(healthCtx, &inventory_v1.GetPartRequest{
		Uuid: gofakeit.UUID(),
	}); err != nil {
		if status.Code(err) == codes.Unavailable {
			return fmt.Errorf("inventory service unavailable: %w", err)
		}
	}

	return nil
}
