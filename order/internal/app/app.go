package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/order/internal/config"
	"github.com/andredubov/rocket-factory/order/internal/metrics"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	httphealth "github.com/andredubov/rocket-factory/platform/pkg/http/health"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	platformMetrics "github.com/andredubov/rocket-factory/platform/pkg/metrics"
	middlewarehttp "github.com/andredubov/rocket-factory/platform/pkg/middleware/http"
	"github.com/andredubov/rocket-factory/platform/pkg/tracing"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  http.Server
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
	// Канал для ошибок от компонентов
	errorsChannel := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errorsChannel <- fmt.Errorf("consumer crashed: %w", err)
		}
	}()

	// HTTP сервер
	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errorsChannel <- fmt.Errorf("http server crashed: %w", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errorsChannel:
		logger.Error(ctx, "❌ component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDIContainer,
		a.initLogger,
		a.initMetrics,
		a.initTracing,
		a.initCloser,
		a.initListener,
		a.initHTTPServer,
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

func (a *App) initMetrics(ctx context.Context) error {
	// Инициализируем платформенный провайдер метрик
	err := platformMetrics.InitProvider(ctx, config.AppConfig().Metrics)
	if err != nil {
		logger.Error(ctx, "❌ failed to init platform metrics provider", zap.Error(err))
		return err
	}

	closer.AddNamed("Metrics provider", func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return platformMetrics.Shutdown(shutdownCtx)
	})

	logger.Info(ctx, "✅ Metrics initialized successfully")

	return metrics.Init(ctx)
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
	listener, err := net.Listen("tcp", config.AppConfig().HTTPServer.Address())
	if err != nil {
		logger.Error(ctx, "❌ failed to init tcp listener", zap.Error(err))
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

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := order_v1.NewServer(a.diContainer.ServerImplementation(ctx))
	if err != nil {
		logger.Error(ctx, "❌ failed to create order server", zap.Error(err))
		return err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/health", httphealth.Handler())

	router.Group(func(r chi.Router) {
		auth := middlewarehttp.NewAuthMiddleware(a.diContainer.AuthClient(ctx))
		r.Use(auth.Handle)
		r.Mount("/", orderServer)
	})

	a.httpServer = http.Server{
		Addr:              a.diContainer.HTTPConfig(ctx).Address(),
		Handler:           router,
		ReadHeaderTimeout: a.diContainer.HTTPConfig(ctx).ReadHeaderTimeout(),
	}

	logger.Info(ctx, "✅ HTTP server initialized successfully")

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	address := config.AppConfig().HTTPServer.Address()
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server starting on %s", address))
	return a.httpServer.Serve(a.listener) // Blocking call
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer running")

	err := a.diContainer.ConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		logger.Error(ctx, "❌ failed to start OrderAssembled Kafka consumer", zap.Error(err))
		return err
	}

	return nil
}
