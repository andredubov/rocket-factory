package app

import (
	"context"
	"log"
	"net"

	"github.com/andredubov/golibs/pkg/closer"
	"github.com/andredubov/golibs/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// App is the main application structure that manages the gRPC server
// and its dependencies through the service provider.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp creates and initializes a new App instance.
// It sets up all required dependencies before returning the application.
func NewApp(ctx context.Context) (*App, error) {
	application := &App{}
	if err := application.initDeps(ctx); err != nil {
		return nil, err
	}

	return application, nil
}

// Run starts the gRPC server and handles graceful shutdown.
// It uses the closer package to ensure proper cleanup of resources.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll() // Close all registered resources
		closer.Wait()     // Wait for cleanup to complete
	}()

	return a.runGRPCServer()
}

// initDeps initializes all application dependencies in sequence.
// Each init function is called and any error will stop the initialization.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,          // Load configuration first
		a.initServiceProvider, // Then setup service provider
		a.initGRPCServer,      // Finally configure gRPC server
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

// initConfig loads the application configuration.
func (a *App) initConfig(_ context.Context) error {
	err := config.Load()
	if err != nil {
		return err
	}

	return nil
}

// initServiceProvider creates a new service provider instance.
// The service provider manages all service dependencies.
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

// initGRPCServer configures and initializes the gRPC server:
// - Uses insecure credentials (for development only)
// - Enables server reflection (for testing)
// - Registers the PaymentService implementation
func (a *App) initGRPCServer(ctx context.Context) error {
	opts := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()), // Disable TLS for development
	}

	a.grpcServer = grpc.NewServer(opts...)
	reflection.Register(a.grpcServer) // Enable reflection API
	payment_v1.RegisterPaymentServiceServer(a.grpcServer, a.serviceProvider.ServerImplementation(ctx))

	return nil
}

// runGRPCServer starts the gRPC server on the configured address.
// It creates a TCP listener and starts serving requests.
func (a *App) runGRPCServer() error {
	addr := a.serviceProvider.GRPCConfig().Address()
	log.Printf("gRPC server starting on %s", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return a.grpcServer.Serve(listener) // Blocking call
}
