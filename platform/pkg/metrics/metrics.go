package metrics

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	defaulTimeout = 5 * time.Second
)

var (
	exporter      *otlpmetricgrpc.Exporter
	meterProvider *metric.MeterProvider
)

type Config interface {
	// CollectorEndpoint возвращает адрес OTLP коллектора
	CollectorEndpoint() string
	// CollectorInterval возвращает интервал отправки метрик
	CollectorInterval() time.Duration
	// CollectorTimeout возвращает таймаут для отправки батча метрик
	CollectorTimeout() time.Duration
}

// InitProvider инициализирует глобальный провайдер метрик OpenTelemetry
func InitProvider(ctx context.Context, cfg Config) error {
	var err error

	// Получаем таймаут из конфигурации или используем по умолчанию
	timeout := cfg.CollectorTimeout()
	if timeout == 0 {
		timeout = defaulTimeout
	}

	// Создание OTLP экспортера с конфигурируемыми параметрами
	exporter, err = otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint()),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithTimeout(timeout),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create metrics exporter")
	}

	// Создание MeterProvider
	meterProvider = metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				exporter,
				metric.WithInterval(cfg.CollectorInterval()),
			),
		),
	)

	otel.SetMeterProvider(meterProvider)

	return nil
}

// SetNoOpMeterProvider устанавливает глобальный NoOp провайдер метрик.
// Полезно для тестов, когда нужно полностью отключить метрики.
func SetNoOpMeterProvider() {
	otel.SetMeterProvider(noop.NewMeterProvider())
}

// GetMeterProvider возвращает текущий провайдер метрик
func GetMeterProvider() *metric.MeterProvider {
	return meterProvider
}

// Shutdown корректно закрывает провайдер метрик и экспортер
func Shutdown(ctx context.Context) error {
	if meterProvider == nil && exporter == nil {
		return nil
	}

	var err error

	if meterProvider != nil {
		err = meterProvider.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown meter provider")
		}
	}

	if exporter != nil {
		err = exporter.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown exporter")
		}
	}

	return nil
}
