package tracer

import (
	"context"
	"errors"
	"fmt"
	"github.com/kholidss/xyz-skilltest/internal/consts"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewExporter(cfg *config.Config) (trace.SpanExporter, error) {
	switch cfg.AppOtelExporter {
	case consts.JaegerExporter:
		return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(fmt.Sprintf("%s:%v/api/traces", cfg.JaegerHost, cfg.JaegerPort))))
	case consts.TempoExporter:
		return tempoExporter(cfg)
	default:
		return nil, errors.New("unknown otel driver")
	}
}

func tempoExporter(cfg *config.Config) (trace.SpanExporter, error) {
	insecureOpt := otlptracehttp.WithInsecure()
	endpointOpt := otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%v", cfg.TempoHost, cfg.TempoPort))

	return otlptracehttp.New(context.Background(), insecureOpt, endpointOpt)
}
