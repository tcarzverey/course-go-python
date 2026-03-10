package telemetry

import (
	"context"
	"fmt"

	// step7 "errors"
	// step7 "time"
	// step7 "go.opentelemetry.io/otel"
	// step7 "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	// step7 "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	// step7 "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// step7 "go.opentelemetry.io/otel/log/global"
	// step7 "go.opentelemetry.io/otel/propagation"
	// step7 "go.opentelemetry.io/otel/sdk/resource"
	// step7 sdklog "go.opentelemetry.io/otel/sdk/log"
	// step7 sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	// step7 sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// step7 semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	// step7 "google.golang.org/grpc"
	// step7 "google.golang.org/grpc/credentials/insecure"
)

// step7 // Init is identical to the url-shortener version — same OTel setup, different service name.
// step7 func Init(ctx context.Context, serviceName, otlpEndpoint string) (func(context.Context) error, error) {
// step7 	conn, err := grpc.NewClient(otlpEndpoint,
// step7 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// step7 	)
// step7 	if err != nil {
// step7 		return nil, fmt.Errorf("grpc connection: %w", err)
// step7 	}
// step7
// step7 	res, err := resource.New(ctx,
// step7 		resource.WithAttributes(semconv.ServiceName(serviceName)),
// step7 	)
// step7 	if err != nil {
// step7 		return nil, fmt.Errorf("resource: %w", err)
// step7 	}
// step7
// step7 	traceExp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
// step7 	if err != nil {
// step7 		return nil, fmt.Errorf("trace exporter: %w", err)
// step7 	}
// step7 	tp := sdktrace.NewTracerProvider(
// step7 		sdktrace.WithBatcher(traceExp),
// step7 		sdktrace.WithResource(res),
// step7 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// step7 	)
// step7 	otel.SetTracerProvider(tp)
// step7 	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
// step7 		propagation.TraceContext{},
// step7 		propagation.Baggage{},
// step7 	))
// step7
// step7 	metricExp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
// step7 	if err != nil {
// step7 		return nil, fmt.Errorf("metric exporter: %w", err)
// step7 	}
// step7 	mp := sdkmetric.NewMeterProvider(
// step7 		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp,
// step7 			sdkmetric.WithInterval(10*time.Second),
// step7 		)),
// step7 		sdkmetric.WithResource(res),
// step7 	)
// step7 	otel.SetMeterProvider(mp)
// step7
// step7 	logExp, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
// step7 	if err != nil {
// step7 		return nil, fmt.Errorf("log exporter: %w", err)
// step7 	}
// step7 	lp := sdklog.NewLoggerProvider(
// step7 		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
// step7 		sdklog.WithResource(res),
// step7 	)
// step7 	global.SetLoggerProvider(lp)
// step7
// step7 	return func(ctx context.Context) error {
// step7 		return errors.Join(tp.Shutdown(ctx), mp.Shutdown(ctx), lp.Shutdown(ctx))
// step7 	}, nil
// step7 }

var _ = fmt.Sprintf
var _ = context.Background
