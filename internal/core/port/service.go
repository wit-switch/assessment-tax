package port

import (
	"context"
	"encoding/csv"

	"github.com/wit-switch/assessment-tax/internal/core/domain"
)

type TaxService interface {
	Calculate(ctx context.Context, body domain.TaxCalculate, allowKReceipt bool) (*domain.Tax, error)
	CalculateFromCSV(ctx context.Context, file csv.Reader) ([]domain.TaxCSV, error)
	UpdateTaxDeduct(ctx context.Context, body domain.UpdateTaxDeduct) (*domain.TaxDeduct, error)
}
