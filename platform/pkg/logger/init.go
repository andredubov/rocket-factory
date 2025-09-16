package logger

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/log"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global       *zap.Logger                // глобальный экземпляр логгера
	initOnce     sync.Once                  // обеспечивает единократную инициализацию
	level        zap.AtomicLevel            // уровень логирования (может изменяться динамически)
	otelProvider *otelLogSdk.LoggerProvider // OTLP provider для graceful shutdown
)

const (
	shutdownTimeout = 2 * time.Second // таймаут для graceful shutdown OTLP provider
)

type Config struct {
	Level              string // уровень логирования ("debug", "info", "warn", "error")
	AsJSON             bool   // формат вывода (true - JSON, false - консольный)
	EnableOTLP         bool   // включение отправки в OpenTelemetry коллектор
	OTLPEndpoint       string // адрес OTLP коллектора
	ServiceName        string // имя сервиса в телеметрии
	ServiceEnvironment string // окружение для фильтрации логов
}

func Init(ctx context.Context, config Config) error {
	initOnce.Do(func() {
		level = zap.NewAtomicLevelAt(parseLevel(config.Level))
		cores := buildCores(ctx, config)
		global = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	})

	if global == nil {
		return fmt.Errorf("zapLogger init failed")
	}

	return nil
}

func buildCores(ctx context.Context, config Config) []zapcore.Core {
	cores := []zapcore.Core{
		createStdoutCore(config.AsJSON),
	}

	if config.EnableOTLP {
		if otlpCore := createOTLPCore(ctx, config.OTLPEndpoint, config.ServiceName, config.ServiceEnvironment); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	return cores
}

// createStdoutCore создает core для записи в stdout/stderr.
// Поддерживает JSON и консольный формат вывода.
func createStdoutCore(asJSON bool) zapcore.Core {
	config := buildEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

// createOTLPCore создает core для отправки в OpenTelemetry коллектор.
// При ошибке подключения возвращает nil (graceful degradation).
func createOTLPCore(ctx context.Context, endpoint, serviceName, serviceEnvironment string) *SimpleOTLPCore {
	otlpLogger, err := createOTLPLogger(ctx, endpoint, serviceName, serviceEnvironment)
	if err != nil {
		return nil
	}

	// Прямо передаём OTLP-логгер в core. Буферизацию делает OTLP SDK (BatchProcessor).
	return NewSimpleOTLPCore(otlpLogger, level)
}

// createOTLPLogger создает OTLP логгер с настроенным экспортером и ресурсами.
// Использует BatchProcessor для эффективной отправки логов.
func createOTLPLogger(ctx context.Context, endpoint, serviceName, serviceEnvironment string) (otelLog.Logger, error) {
	exporter, err := createOTLPExporter(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	rs, err := createResource(ctx, serviceName, serviceEnvironment)
	if err != nil {
		return nil, err
	}

	provider := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)
	otelProvider = provider

	return provider.Logger("app"), nil
}

// createOTLPExporter создает gRPC экспортер для OTLP коллектора
func createOTLPExporter(ctx context.Context, endpoint string) (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(), // для разработки, в продакшене следует использовать TLS
	)
}

// createResource создает метаданные сервиса для телеметрии
func createResource(ctx context.Context, serviceName, serviceEnvironment string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment", serviceEnvironment),
		),
	)
}

// buildEncoderConfig настраивает формат вывода логов с нужными полями
func buildEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}

// parseLevel преобразует строковое значение в zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
