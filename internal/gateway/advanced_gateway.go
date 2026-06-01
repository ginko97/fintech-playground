package gateway

import (
	"context"
	"time"

	"github.com/ginko97/fintech-playground/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AdvancedGateway struct {
	tracer trace.Tracer
}

func NewAdvancedGateway() *AdvancedGateway {
	return &AdvancedGateway{
		tracer: otel.Tracer("advanced-gateway"),
	}
}

// ProcessPayment calls external PSP with tracing
func (g *AdvancedGateway) ProcessPayment(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	ctx, span := g.tracer.Start(ctx, "ProcessPayment")
	defer span.End()

	// Simulate network call to external PSP
	time.Sleep(300 * time.Millisecond)

	// Simulate business logic
	if tx.Amount%7 == 0 {
		tx.Status = domain.StatusFailed
		span.SetAttributes(
			attribute.String("payment.result", "failed"),
			attribute.Int64("payment.amount", tx.Amount),
		)
	} else {
		tx.Status = domain.StatusCompleted
		span.SetAttributes(
			attribute.String("payment.result", "success"),
			attribute.Int64("payment.amount", tx.Amount),
		)
	}

	return tx, nil
}
