package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/andredubov/rocket-factory/payment/internal/config"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	"github.com/andredubov/rocket-factory/platform/pkg/grpc/health"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// App represents the main application structure with dependencies.
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

// New creates and initializes a new App instance with all dependencies.
func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run starts the gRPC server and begins serving requests.
func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

// initDeps initializes all application dependencies in sequence.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDIContainer,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
	}

	for i, f := range inits {
		err := f(ctx)
		if err != nil {
			return fmt.Errorf("init step %d failed: %w", i, err)
		}
	}

	return nil
}

// initDIContainer initializes the dependency injection container.
func (a *App) initDIContainer(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

// initLogger initializes the application logger with configured settings.
func (a *App) initLogger(ctx context.Context) error {
	loggerConfig := logger.Config{
		Level:              config.AppConfig().Logger.Level(),
		AsJSON:             config.AppConfig().Logger.AsJson(),
		EnableOTLP:         config.AppConfig().Logger.EnableOTLP(),
		OTLPEndpoint:       config.AppConfig().Logger.OTLPEndpoint(),
		ServiceName:        config.AppConfig().Logger.ServiceName(),
		ServiceEnvironment: config.AppConfig().Logger.ServiceEnvironment(),
	}
	return logger.Init(ctx, loggerConfig)
}

// initCloser sets up the application closer with the configured logger.
func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

// initListener creates a TCP listener on the configured gRPC server address.
func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().GRPCServer.Address())
	if err != nil {
		return err
	}

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

// initGRPCServer configures and initializes the gRPC server with required services.
func (a *App) initGRPCServer(ctx context.Context) error {
	authInterceptor := a.diContainer.AuthInterceptor(ctx)
	if authInterceptor == nil {
		return errors.New("failed to initialize auth interceptor")
	}

	opts := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()), // Disable TLS for development
		grpc.ChainUnaryInterceptor(
			authInterceptor.Unary(),
		),
	}

	a.grpcServer = grpc.NewServer(opts...)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer) // Enable reflection API

	health.RegisterService(a.grpcServer) // for healthcheck

	payment_v1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.ServerImplementation(ctx))

	return nil
}

// runGRPCServer starts the gRPC server and begins listening for incoming requests.
func (a *App) runGRPCServer(ctx context.Context) error {
	address := config.AppConfig().GRPCServer.Address()
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService server starting on %s", address))

	return a.grpcServer.Serve(a.listener) // Blocking call
}
