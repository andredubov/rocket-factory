package app

import (
	"context"
	"log"

	"github.com/andredubov/golibs/pkg/config"
	"github.com/andredubov/golibs/pkg/config/env"

	server "github.com/andredubov/rocket-factory/payment/internal/api/v1/payment"
	"github.com/andredubov/rocket-factory/payment/internal/service"
	"github.com/andredubov/rocket-factory/payment/internal/service/payment"
)

// serviceProvider implements the dependency injection container pattern.
// It lazily initializes and provides access to all service dependencies.
type serviceProvider struct {
	paymentService       service.Payments              // Business logic service
	grpcConfig           config.GRPCConfig             // gRPC server configuration
	serverImplementation *server.PaymentImplementation // gRPC handler implementation
}

// newServiceProvider creates a new service provider instance.
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// GRPCConfig loads and provides the gRPC server configuration.
// Implements lazy initialization - config is loaded only once.
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
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
func (s *serviceProvider) PaymentService(ctx context.Context) service.Payments {
	if s.paymentService == nil {
		s.paymentService = payment.NewService()
	}

	return s.paymentService
}

// ServerImplementation creates and provides the gRPC server implementation.
// It initializes all required dependencies (payment service) automatically.
func (s *serviceProvider) ServerImplementation(ctx context.Context) *server.PaymentImplementation {
	if s.serverImplementation == nil {
		paymentService := s.PaymentService(ctx)
		s.serverImplementation = server.NewPaymentImplementation(paymentService)
	}

	return s.serverImplementation
}
