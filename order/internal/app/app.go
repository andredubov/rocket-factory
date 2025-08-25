package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/order/internal/config"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
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
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
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
		a.initCloser,
		a.initListener,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDIContainer(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().HTTPServer.Address())
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

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := order_v1.NewServer(a.diContainer.ServerImplementation(ctx))
	if err != nil {
		log.Printf("failed to create order server: %v\n", err)
		return err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/", orderServer)

	a.httpServer = http.Server{
		Addr:              a.diContainer.HTTPConfig().Address(),
		Handler:           router,
		ReadHeaderTimeout: a.diContainer.HTTPConfig().ReadHeaderTimeout(),
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	address := config.AppConfig().HTTPServer.Address()
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC OrderService server starting on %s", address))

	return a.httpServer.Serve(a.listener) // Blocking call
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer running")

	err := a.diContainer.ConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
