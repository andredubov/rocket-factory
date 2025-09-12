package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName                 = "assembly-service"
	assemblyDurationMetricName  = "assembly_duration_seconds"
	assemblyDurationDescription = "Время сборки корабля в секундах"
	secondsUnit                 = "s"
)

// assemblyDurationSeconds - гистограмма длительности сборки
var assemblyDurationSeconds metric.Float64Histogram

// Init инициализирует все инструменты метрик с описаниями для production
func Init(_ context.Context) error {
	var err error

	meter := otel.Meter(serviceName)

	// Создание гистограммы времени сборки
	assemblyDurationSeconds, err = meter.Float64Histogram(
		assemblyDurationMetricName,
		metric.WithDescription(assemblyDurationDescription),
		metric.WithUnit(secondsUnit),
		metric.WithExplicitBucketBoundaries(
			// Бакеты от 1 секунды до 30 секунд
			1.0, 2.0, 5.0, 8.0, 10.0, 12.0, 15.0, 20.0, 25.0, 30.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}

// RecordAssemblyDuration записывает время сборки
func RecordAssemblyDuration(ctx context.Context, duration time.Duration) {
	assemblyDurationSeconds.Record(ctx, duration.Seconds())
}
