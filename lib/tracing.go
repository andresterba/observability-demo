package lib

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type TraceExportTarget int

const (
	Stdout TraceExportTarget = iota
	Backend
)

func GetTracer(ctx context.Context, target TraceExportTarget) (*sdktrace.TracerProvider, error) {
	var exporter sdktrace.SpanExporter
	var err error

	switch target {
	case Stdout:
		exporter, err = stdout.New(stdout.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
	case Backend:
		exporter, err = otlptracehttp.New(
			ctx,
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint("127.0.0.1:4318"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}

	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}

func SetRuntimeSettings(serviceName string) {
	_ = os.Setenv("OTEL_SERVICE_NAME", serviceName)
}
