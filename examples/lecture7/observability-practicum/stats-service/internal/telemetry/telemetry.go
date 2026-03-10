package telemetry

import (
	"context"
	"fmt"

	// step7 "go.opentelemetry.io/otel"
	// step7 "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// step7 "go.opentelemetry.io/otel/propagation"
	// step7 "go.opentelemetry.io/otel/sdk/resource"
	// step7 sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// step7 semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	// step7 "google.golang.org/grpc"
	// step7 "google.golang.org/grpc/credentials/insecure"
)

// step7 // InitTracer is identical to the url-shortener version — same OTel setup, different service name.
// step7 func InitTracer(ctx context.Context, serviceName, otlpEndpoint string) (func(context.Context) error, error) {
// step7 	conn, err := grpc.NewClient(otlpEndpoint,
// step7 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// step7 	)
// step7 	if err != nil {
// step7 		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
// step7 	}
// step7
// step7 	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
// step7 	if err != nil {
// step7 		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
// step7 	}
// step7
// step7 	res, err := resource.New(ctx,
// step7 		resource.WithAttributes(semconv.ServiceName(serviceName)),
// step7 	)
// step7 	if err != nil {
// step7 		return nil, fmt.Errorf("failed to create resource: %w", err)
// step7 	}
// step7
// step7 	tp := sdktrace.NewTracerProvider(
// step7 		sdktrace.WithBatcher(exporter),
// step7 		sdktrace.WithResource(res),
// step7 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// step7 	)
// step7
// step7 	otel.SetTracerProvider(tp)
// step7 	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
// step7 		propagation.TraceContext{},
// step7 		propagation.Baggage{},
// step7 	))
// step7
// step7 	return tp.Shutdown, nil
// step7 }

var _ = fmt.Sprintf
var _ = context.Background
