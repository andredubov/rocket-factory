package logger

import (
	"context"

	"go.uber.org/zap"
)

type logger struct {
	zapLogger *zap.Logger
}

// Debug записывает лог уровня DEBUG с контекстом.
func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if l.zapLogger != nil {
		l.zapLogger.Debug(msg, fields...)
	}
}

// Info записывает лог уровня INFO с контекстом.
func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if l.zapLogger != nil {
		l.zapLogger.Info(msg, fields...)
	}
}

// Warn записывает лог уровня WARN с контекстом.
func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if l.zapLogger != nil {
		l.zapLogger.Warn(msg, fields...)
	}
}

// Error записывает лог уровня ERROR с контекстом.
func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if l.zapLogger != nil {
		l.zapLogger.Error(msg, fields...)
	}
}

// Fatal записывает лог уровня FATAL с контекстом и завершает программу.
func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if l.zapLogger != nil {
		l.zapLogger.Fatal(msg, fields...)
	}
}

// Sync принудительно сбрасывает все буферизованные логи.
// Должен вызываться перед завершением работы приложения.
func (l *logger) Sync() error {
	if l.zapLogger != nil {
		return l.zapLogger.Sync()
	}
	return nil
}

// Close корректно завершает работу логгера.
// Вызывает Sync и освобождает ресурсы.
func (l *logger) Close() error {
	var syncErr error
	if l.zapLogger != nil {
		syncErr = l.zapLogger.Sync()
	}

	// Дополнительно вызываем глобальный Close для OTLP provider
	closeErr := Close()

	// Возвращаем первую ошибку, если есть
	if syncErr != nil {
		return syncErr
	}
	return closeErr
}

func Logger() *logger {
	return &logger{zapLogger: global}
}
