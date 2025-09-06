package app

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/andredubov/rocket-factory/payment/internal/api/v1/payment"
	"github.com/andredubov/rocket-factory/payment/internal/config"
	"github.com/andredubov/rocket-factory/payment/internal/config/env"
	"github.com/andredubov/rocket-factory/payment/internal/service/payment"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	middlewaregrpc "github.com/andredubov/rocket-factory/platform/pkg/middleware/grpc"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
)

// diContainer  implements the dependency injection container pattern.
type diContainer struct {
	paymentService       api.PaymentService         // Business logic service
	grpcConfig           config.GRPCConfig          // gRPC server configuration
	serverImplementation *api.PaymentImplementation // gRPC handler implementation

	authClient      auth_v1.AuthServiceClient
	authInterceptor *middlewaregrpc.AuthInterceptor
}

// newDIContainer creates a new service provider instance.
func NewDIContainer() *diContainer {
	return &diContainer{}
}

// GRPCConfig loads and provides the gRPC server configuration.
func (s *diContainer) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// PaymentService provides the payment business logic service.
func (s *diContainer) PaymentService(ctx context.Context) api.PaymentService {
	if s.paymentService == nil {
		s.paymentService = payment.NewService()
	}

	return s.paymentService
}

// ServerImplementation creates and provides the gRPC server implementation.
func (s *diContainer) ServerImplementation(ctx context.Context) *api.PaymentImplementation {
	if s.serverImplementation == nil {
		paymentService := s.PaymentService(ctx)
		s.serverImplementation = api.NewPaymentImplementation(paymentService)
	}

	return s.serverImplementation
}

// AuthClient creates and provides the IAM-service grpc client
func (d *diContainer) AuthClient(ctx context.Context) auth_v1.AuthServiceClient {
	if d.authClient == nil {
		iamAddress := config.AppConfig().IAMClient.Address()
		if iamAddress == "" {
			logger.Error(ctx, "IAM client address is empty")
			return nil
		}

		connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		iamConn, err := grpc.NewClient(
			iamAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(connectCtx, "Failed to connect to IAM service", zap.Error(err))
			return nil
		}

		closer.AddNamed("IAM grpc connection", func(ctx context.Context) error {
			return iamConn.Close()
		})

		d.authClient = auth_v1.NewAuthServiceClient(iamConn)
	}

	return d.authClient
}

// AuthInterceptor creates and provides the auth interceptor
func (d *diContainer) AuthInterceptor(ctx context.Context) *middlewaregrpc.AuthInterceptor {
	if d.authInterceptor == nil {
		d.authInterceptor = middlewaregrpc.NewAuthInterceptor(
			d.AuthClient(ctx),
		)
	}

	return d.authInterceptor
}
