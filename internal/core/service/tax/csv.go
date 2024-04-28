package tax

import (
	"context"
	"encoding/csv"

	"github.com/gocarina/gocsv"
	"github.com/shopspring/decimal"
	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/pkg/errorx"
)

func (s *Services) CalculateFromCSV(ctx context.Context, file csv.Reader) ([]domain.TaxCSV, error) {
	var csv []domain.CSVFile
	if err := gocsv.UnmarshalCSV(&file, &csv); err != nil {
		return nil, err
	}

	var zero decimal.Decimal
	out := make([]domain.TaxCSV, len(csv))
	for i, v := range csv {
		if v.TotalIncome.Decimal.Cmp(zero) < 0 {
			params := domain.NewFieldMessageList().Add(
				"totalIncome", "totalIncome should greater than or equal 0",
			).Value()
			return nil, errorx.ErrValidationFail.WithParams(params)
		}

		if v.Wht.Decimal.Cmp(zero) < 0 {
			params := domain.NewFieldMessageList().Add(
				"wht", "wht should greater than or equal 0",
			).Value()
			return nil, errorx.ErrValidationFail.WithParams(params)
		}

		if v.Wht.Decimal.Cmp(v.TotalIncome.Decimal) > 0 {
			params := domain.NewFieldMessageList().Add(
				"wht", "wht should less than or equal totalIncome",
			).Value()
			return nil, errorx.ErrValidationFail.WithParams(params)
		}

		taxCal, err := s.Calculate(ctx, domain.TaxCalculate{
			TotalIncome: v.TotalIncome.Decimal,
			Wht:         v.Wht.Decimal,
			Allowances: []domain.Allowance{{
				AllowanceType: domain.TaxDeductTypeDonation,
				Amount:        v.Donation.Decimal,
			}},
		}, false)
		if err != nil {
			return nil, err
		}

		out[i] = domain.TaxCSV{
			TotalIncome: v.TotalIncome.Decimal,
			Tax:         taxCal.Tax,
			TaxRefund:   taxCal.TaxRefund,
		}
	}

	return out, nil
}
