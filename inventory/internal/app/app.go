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

	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// App represents the core application structure
// Manages lifecycle and dependencies of the service
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp constructs new application instance
// Initializes all dependencies through initDeps
func NewApp(ctx context.Context) (*App, error) {
	application := &App{}
	if err := application.initDeps(ctx); err != nil {
		return nil, err
	}

	return application, nil
}

// Run starts the application:
// 1. Launches GRPC server
// 2. Sets up graceful shutdown via closer
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

// initDeps initializes application dependencies
// Executes initialization functions sequentially
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,          // Load configuration
		a.initServiceProvider, // Initialize service container
		a.initGRPCServer,      // Setup GRPC server
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

// initConfig loads application configuration
func (a *App) initConfig(_ context.Context) error {
	err := config.Load()
	if err != nil {
		return err
	}

	return nil
}

// initServiceProvider creates service provider instance
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

// initGRPCServer configures GRPC server:
// 1. Creates server with insecure credentials (dev only)
// 2. Enables reflection for testing
// 3. Registers inventory service
func (a *App) initGRPCServer(ctx context.Context) error {
	opts := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()), // Disabled security (development only)
	}

	a.grpcServer = grpc.NewServer(opts...)
	reflection.Register(a.grpcServer) // For grpcurl testing
	inventory_v1.RegisterInventoryServiceServer(a.grpcServer, a.serviceProvider.ServerImplementation(ctx))

	return nil
}

// runGRPCServer starts GRPC server on configured address
func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	listener, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	if err = a.grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}
