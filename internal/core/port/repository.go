package port

import (
	"context"

	"github.com/wit-switch/assessment-tax/internal/core/domain"
)

type TaxRepository interface {
	GetTaxDeductByType(ctx context.Context, q domain.TaxDeductType) (*domain.TaxDeduct, error)
	ListTaxDeduct(ctx context.Context, q domain.GetTaxDeduct) ([]domain.TaxDeduct, error)
	UpdateTaxDeduct(ctx context.Context, u domain.UpdateTaxDeduct) (*domain.TaxDeduct, error)
}
