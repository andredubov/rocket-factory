package app

import (
	"context"
	"log"

	api "github.com/andredubov/rocket-factory/payment/internal/api/v1/payment"
	"github.com/andredubov/rocket-factory/payment/internal/config"
	"github.com/andredubov/rocket-factory/payment/internal/config/env"
	"github.com/andredubov/rocket-factory/payment/internal/service/payment"
)

// diContainer  implements the dependency injection container pattern.
// It lazily initializes and provides access to all service dependencies.
type diContainer struct {
	paymentService       api.PaymentService         // Business logic service
	grpcConfig           config.GRPCConfig          // gRPC server configuration
	serverImplementation *api.PaymentImplementation // gRPC handler implementation
}

// newDIContainer creates a new service provider instance.
func NewDIContainer() *diContainer {
	return &diContainer{}
}

// GRPCConfig loads and provides the gRPC server configuration.
// Implements lazy initialization - config is loaded only once.
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
// Initializes the service only when first requested.
func (s *diContainer) PaymentService(ctx context.Context) api.PaymentService {
	if s.paymentService == nil {
		s.paymentService = payment.NewService()
	}

	return s.paymentService
}

// ServerImplementation creates and provides the gRPC server implementation.
// It initializes all required dependencies (payment service) automatically.
func (s *diContainer) ServerImplementation(ctx context.Context) *api.PaymentImplementation {
	if s.serverImplementation == nil {
		paymentService := s.PaymentService(ctx)
		s.serverImplementation = api.NewPaymentImplementation(paymentService)
	}

	return s.serverImplementation
}
