package app

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/notification/internal/config"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
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

	go func() {
		if err := a.runOrderPaidEventConsumer(ctx); err != nil {
			errorsChannel <- fmt.Errorf("order paid event consumer crashed: %w", err)
		}
	}()

	go func() {
		if err := a.runOrderAssembledEventConsumer(ctx); err != nil {
			errorsChannel <- fmt.Errorf("order assembled event consumer crashed: %w", err)
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
		a.initTelegramBot,
	}

	for i, f := range inits {
		err := f(ctx)
		if err != nil {
			return fmt.Errorf("init step %d failed: %w", i, err)
		}
	}

	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	// Получаем бота из DI контейнера
	telegramBot := a.diContainer.TelegramBot(ctx)

	// Регистрируем обработчик для активации бота
	telegramBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		logger.Info(ctx, "chat id", zap.Int64("chat_id", update.Message.Chat.ID))

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "🔔 Notification Bot активирован! Теперь вы будете получать уведомления о оплате частей и сборки корабля.",
		})
		if err != nil {
			logger.Error(ctx, "❌ failed to send activation message", zap.Error(err))
		}
	})

	// Запускаем бота в фоне
	go func() {
		logger.Info(ctx, "🚀 Telegram bot started...")
		telegramBot.Start(ctx)
	}()

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

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runOrderPaidEventConsumer(ctx context.Context) error {
	err := a.diContainer.OrderPaidEventConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		logger.Error(ctx, "❌ failed to start OrderPaidEvent Kafka consumer", zap.Error(err))
		return err
	}

	logger.Info(ctx, "🚀 OrderPaidEvent Kafka consumer successfully started and running")

	return nil
}

func (a *App) runOrderAssembledEventConsumer(ctx context.Context) error {
	err := a.diContainer.OrderAssembledEventConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Failed to start OrderAssembledEvent Kafka consumer", zap.Error(err))
		return err
	}

	logger.Info(ctx, "🚀 OrderAssembledEvent Kafka consumer successfully started and running")

	return nil
}
