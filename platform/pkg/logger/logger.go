package logger

import (
	"context"

	"go.uber.org/zap"
)

// Info записывает лог уровня INFO.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Info(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Info(msg, fields...)
	}
}

// Info записывает лог уровня Debug.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Debug(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Debug(msg, fields...)
	}
}

// Warn записывает лог уровня WARN.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Warn(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Warn(msg, fields...)
	}
}

// Error записывает лог уровня ERROR.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Error(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Error(msg, fields...)
	}
}

// Fatal записывает лог уровня FATAL и завершает программу.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Fatal(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Fatal(msg, fields...)
	}
}

// SetNopLogger устанавливает nop логгер (для тестов)
func SetNopLogger() {
	global = zap.NewNop()
}

// SetTestLogger устанавливает тестовый логгер для использования в тестах
func SetTestLogger(testLogger *zap.Logger) {
	global = testLogger
}

// GetGlobalLogger возвращает глобальный логгер (для тестов)
func GetGlobalLogger() *zap.Logger {
	return global
}

// Sync принудительно сбрасывает все буферизованные логи.
// Вызывает sync для всех cores (stdout + OTLP).
func Sync() error {
	if global != nil {
		return global.Sync()
	}

	return nil
}

// Close корректно завершает работу логгера.
// Останавливает OTLP provider с таймаутом для отправки оставшихся логов.
func Close() error {
	if otelProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := otelProvider.Shutdown(ctx); err != nil {
			if global != nil {
				global.Error("failed to shutdown OTLP provider", zap.Error(err))
			}
		}
	}

	return nil
}
