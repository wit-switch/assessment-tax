package port

import (
	"context"

	"github.com/wit-switch/assessment-tax/internal/core/domain"
)

type TaxService interface {
	Calculate(ctx context.Context, body domain.TaxCalculate) (*domain.Tax, error)
}
