package app

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/assembly/internal/config"
	"github.com/andredubov/rocket-factory/assembly/internal/metrics"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	platformMetrics "github.com/andredubov/rocket-factory/platform/pkg/metrics"
)

type App struct {
	diContainer *diContainer
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
		a.initCloser,
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
	if err := platformMetrics.InitProvider(ctx, config.AppConfig().Metrics); err != nil {
		logger.Error(ctx, "Failed to init platform metrics provider", zap.Error(err))
		return err
	}

	closer.AddNamed("Metrics provider", func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return platformMetrics.Shutdown(shutdownCtx)
	})

	return metrics.Init(ctx)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	err := a.diContainer.ConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		logger.Error(ctx, "❌ failed to start OrderPaid Kafka consumer", zap.Error(err))
		return err
	}

	logger.Info(ctx, "🚀 OrderPaid Kafka consumer running")

	return nil
}
