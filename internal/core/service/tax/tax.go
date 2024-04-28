package tax

import (
	"context"

	"github.com/wit-switch/assessment-tax/internal/core/domain"

	"github.com/shopspring/decimal"
)

func (s *Services) Calculate(ctx context.Context, body domain.TaxCalculate) (*domain.Tax, error) {
	taxDeducts, err := s.taxRepository.ListTaxDeduct(ctx, domain.GetTaxDeduct{})
	if err != nil {
		return nil, err
	}

	taxDeductMap := make(map[domain.TaxDeductType]domain.TaxDeduct)
	for _, v := range taxDeducts {
		taxDeductMap[v.Type] = v
	}

	personalDeduct, err := GetTaxDeductByType(taxDeductMap, domain.TaxDeductTypePersonal)
	if err != nil {
		return nil, err
	}

	personal := personalDeduct.Amount

	totalIncome := body.TotalIncome.Sub(personal)

	var totalTax decimal.Decimal
	for _, taxRate := range taxRates {
		if totalIncome.Cmp(taxRate.MinIncome) > 0 {
			diff := totalIncome.Sub(taxRate.MinIncome)
			totalIncome = totalIncome.Sub(diff)
			totalTax = totalTax.Add(diff.Mul(taxRate.Rate))
		}
	}

	return &domain.Tax{
		Tax: totalTax,
	}, nil
}
