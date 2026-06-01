package infrastructure

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

var Tracer = otel.Tracer("fintech-playground")

func InitTracing() {
	// Simple stdout exporter first (easier to debug)
	exporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
}

func ShutdownTracing() {
}
