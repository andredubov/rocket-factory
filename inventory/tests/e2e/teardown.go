//go:build integration

package integration

import (
	"context"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

// teardownTestEnvironment — освобождает все ресурсы тестового окружения
func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "🧹 Очистка тестового окружения...")

	cleanupTestEnvironment(ctx, env)

	log.Info(ctx, "✅ Тестовое окружение успешно очищено")
}

// cleanupTestEnvironment — вспомогательная функция для освобождения ресурсов
func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.InventoryAppContainer != nil {
		if err := env.InventoryAppContainer.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер Inventory сервиса", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер Inventory-сервиса остановлен")
		}
	}

	if env.InventoryMongoContainer != nil {
		if err := env.InventoryMongoContainer.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер MongoDB", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер MongoDB остановлен")
		}
	}

	if env.IAMAppContainer != nil {
		if err := env.IAMAppContainer.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер IAM сервиса", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер IAM сервиса остановлен")
		}
	}

	if env.IAMPostgresContainer != nil {
		if err := env.IAMPostgresContainer.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер PostgresDB", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер PostgresDB остановлен")
		}
	}

	if env.IAMRedisContainer != nil {
		if err := env.IAMRedisContainer.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер Redis", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер Redis остановлен")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "не удалось удалить сеть", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Сеть удалена")
		}
	}
}
