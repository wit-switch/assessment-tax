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

	donationDeduct, err := GetTaxDeductByType(taxDeductMap, domain.TaxDeductTypeDonation)
	if err != nil {
		return nil, err
	}

	allowanceMap := make(map[domain.TaxDeductType]decimal.Decimal)
	for _, v := range body.Allowances {
		allowanceMap[v.AllowanceType] = v.Amount
	}

	personal := personalDeduct.Amount
	donation, err := GetAllowanceAmount(ctx, donationDeduct, allowanceMap[domain.TaxDeductTypeDonation])
	if err != nil {
		return nil, err
	}

	totalIncome := body.TotalIncome.Sub(personal).Sub(donation)

	var totalTax, taxRefund, zero decimal.Decimal
	taxLevel := make([]domain.TaxLevel, 0, len(taxRates))
	for _, taxRate := range taxRates {
		level := domain.TaxLevel{
			Level: taxRate.Description,
			Tax:   zero,
		}
		if totalIncome.Cmp(taxRate.MinIncome) > 0 {
			diff := totalIncome.Sub(taxRate.MinIncome)
			calTax := diff.Mul(taxRate.Rate)
			totalIncome = totalIncome.Sub(diff)
			totalTax = totalTax.Add(calTax)
			level.Tax = calTax
		}

		taxLevel = append([]domain.TaxLevel{level}, taxLevel...)
	}

	tax := totalTax.Sub(body.Wht)
	if tax.Cmp(zero) < 0 {
		taxRefund = tax.Abs()
		tax = zero
	}

	return &domain.Tax{
		Tax:       tax,
		TaxRefund: taxRefund,
		TaxLevel:  taxLevel,
	}, nil
}
