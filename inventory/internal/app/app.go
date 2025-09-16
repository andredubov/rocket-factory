package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/andredubov/rocket-factory/inventory/internal/config"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	"github.com/andredubov/rocket-factory/platform/pkg/grpc/health"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	"github.com/andredubov/rocket-factory/platform/pkg/tracing"
	"github.com/andredubov/rocket-factory/shared/pkg/interceptors"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDIContainer,
		a.initLogger,
		a.initTracing,
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

func (a *App) initDIContainer(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

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

func (a *App) initTracing(ctx context.Context) error {
	err := tracing.InitTracer(ctx, config.AppConfig().Tracing)
	if err != nil {
		logger.Error(ctx, "❌ failed to init tracing", zap.Error(err))
		return err
	}

	closer.AddNamed("Tracing", func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return tracing.ShutdownTracer(shutdownCtx)
	})

	logger.Info(ctx, "✅ Tracing initialized successfully")
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(ctx context.Context) error {
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

	logger.Info(ctx, "✅ TCP listener initialized successfully")

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	authInterceptor := a.diContainer.AuthInterceptor(ctx)
	if authInterceptor == nil {
		return errors.New("failed to initialize auth interceptor")
	}

	opts := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()), // Disable TLS for development
		grpc.ChainUnaryInterceptor(
			interceptors.UnaryErrorInterceptor(),
			authInterceptor.Unary(),
			tracing.UnaryServerInterceptor(config.AppConfig().Tracing.ServiceName()),
		),
	}

	a.grpcServer = grpc.NewServer(opts...)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)    // Enable reflection API
	health.RegisterService(a.grpcServer) // for healthcheck
	inventory_v1.RegisterInventoryServiceServer(a.grpcServer, a.diContainer.ServerImplementation(ctx))

	logger.Info(ctx, "✅ gRPC server initialized successfully")

	return nil
}

// runGRPCServer starts the gRPC server and begins listening for incoming requests.
func (a *App) runGRPCServer(ctx context.Context) error {
	address := config.AppConfig().GRPCServer.Address()
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC InventoryService server starting on %s", address))

	return a.grpcServer.Serve(a.listener) // Blocking call
}
