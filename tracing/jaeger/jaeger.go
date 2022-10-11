package jaeger

import (
	"context"

	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/tracing"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	// CollectorEndpoint is the endpoint for jaeger span collection.
	CollectorEndpoint string `mapstructure:"collectorEndpoint"`
}

// Jaeger is the NamedLoader for Jaeger.
func Jaeger() (string, tracing.Loader) {
	return "jaeger", Loader
}

func Loader(ctx context.Context, with interface{}, resolveAs resolve.ResolveAs) (trace.SpanExporter, error) {
	c := Config{
		CollectorEndpoint: "http://localhost:14268/api/traces",
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(c.CollectorEndpoint)))
}
