package telemetry

import (
	"context"
	"fmt"

	// step2 "errors"
	// step2 "time"
	// step2 "go.opentelemetry.io/otel"
	// step2 "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	// step2 "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	// step2 "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// step2 "go.opentelemetry.io/otel/log/global"
	// step2 "go.opentelemetry.io/otel/propagation"
	// step2 "go.opentelemetry.io/otel/sdk/resource"
	// step2 sdklog "go.opentelemetry.io/otel/sdk/log"
	// step2 sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	// step2 sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// step2 semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	// step2 "google.golang.org/grpc"
	// step2 "google.golang.org/grpc/credentials/insecure"
)

// step2 // Init initialises all three OTel signal providers over a shared gRPC connection:
// step2 //   - TracerProvider  → traces → OTel Collector → Tempo / Jaeger
// step2 //   - MeterProvider   → metrics → OTel Collector → Prometheus
// step2 //   - LoggerProvider  → logs → OTel Collector → Loki
// step2 //
// step2 // Each provider is registered as the global default, so otelhttp, otelslog and
// step2 // otel.Tracer() all pick it up automatically without extra wiring.
// step2 func Init(ctx context.Context, serviceName, otlpEndpoint string) (func(context.Context) error, error) {
// step2 	conn, err := grpc.NewClient(otlpEndpoint,
// step2 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// step2 	)
// step2 	if err != nil {
// step2 		return nil, fmt.Errorf("grpc connection: %w", err)
// step2 	}
// step2
// step2 	res, err := resource.New(ctx,
// step2 		resource.WithAttributes(semconv.ServiceName(serviceName)),
// step2 	)
// step2 	if err != nil {
// step2 		return nil, fmt.Errorf("resource: %w", err)
// step2 	}
// step2
// step2 	// ── Traces ───────────────────────────────────────────────────────────────
// step2 	traceExp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
// step2 	if err != nil {
// step2 		return nil, fmt.Errorf("trace exporter: %w", err)
// step2 	}
// step2 	tp := sdktrace.NewTracerProvider(
// step2 		sdktrace.WithBatcher(traceExp),
// step2 		sdktrace.WithResource(res),
// step2 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// step2 	)
// step2 	otel.SetTracerProvider(tp)
// step2 	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
// step2 		propagation.TraceContext{},
// step2 		propagation.Baggage{},
// step2 	))
// step2
// step2 	// ── Metrics ──────────────────────────────────────────────────────────────
// step2 	metricExp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
// step2 	if err != nil {
// step2 		return nil, fmt.Errorf("metric exporter: %w", err)
// step2 	}
// step2 	mp := sdkmetric.NewMeterProvider(
// step2 		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp,
// step2 			sdkmetric.WithInterval(10*time.Second),
// step2 		)),
// step2 		sdkmetric.WithResource(res),
// step2 	)
// step2 	otel.SetMeterProvider(mp)
// step2
// step2 	// ── Logs ─────────────────────────────────────────────────────────────────
// step2 	logExp, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
// step2 	if err != nil {
// step2 		return nil, fmt.Errorf("log exporter: %w", err)
// step2 	}
// step2 	lp := sdklog.NewLoggerProvider(
// step2 		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
// step2 		sdklog.WithResource(res),
// step2 	)
// step2 	global.SetLoggerProvider(lp)
// step2
// step2 	return func(ctx context.Context) error {
// step2 		return errors.Join(tp.Shutdown(ctx), mp.Shutdown(ctx), lp.Shutdown(ctx))
// step2 	}, nil
// step2 }

// placeholder so the package compiles without step2
var _ = fmt.Sprintf
var _ = context.Background
