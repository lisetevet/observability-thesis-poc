package observability

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// InitProvider configures a TracerProvider and global propagator.
// It relies on standard OTEL env vars (OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_EXPORTER_OTLP_HEADERS, etc).
func InitProvider(ctx context.Context, serviceName string) func(context.Context) error {
	// Convenience for New Relic: allow using NEW_RELIC_LICENSE_KEY without manually setting headers.
	if os.Getenv("OTEL_EXPORTER_OTLP_HEADERS") == "" {
		if key := os.Getenv("NEW_RELIC_LICENSE_KEY"); key != "" {
			_ = os.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "api-key="+key)
		}
	}

	exp, err := otlptracehttp.New(ctx) // reads OTEL_* env vars
	if err != nil {
		// If exporter init fails, keep app running (but without traces).
		// You can change this to panic/log.Fatal later if you prefer.
		return func(context.Context) error { return nil }
	}

	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown
}
