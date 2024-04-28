package port

import (
	"context"

	"github.com/wit-switch/assessment-tax/internal/core/domain"
)

type TaxRepository interface {
	ListTaxDeduct(ctx context.Context, q domain.GetTaxDeduct) ([]domain.TaxDeduct, error)
}
