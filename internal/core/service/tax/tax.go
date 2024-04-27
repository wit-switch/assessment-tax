package tax

import (
	"context"
	"fmt"

	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/pkg/errorx"

	"github.com/shopspring/decimal"
)

func (s *Services) Calculate(ctx context.Context, body domain.TaxCalculate, allowKReceipt bool) (*domain.Tax, error) {
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

	var kReceipt decimal.Decimal
	if allowKReceipt {
		kReceiptDeduct, err := GetTaxDeductByType(taxDeductMap, domain.TaxDeductTypeKReceipt)
		if err != nil {
			return nil, err
		}
		kReceipt, err = GetAllowanceAmount(ctx, kReceiptDeduct, allowanceMap[domain.TaxDeductTypeKReceipt])
		if err != nil {
			return nil, err
		}
	}

	totalIncome := body.TotalIncome.Sub(personal).Sub(donation).Sub(kReceipt)

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

func (s *Services) UpdateTaxDeduct(ctx context.Context, body domain.UpdateTaxDeduct) (*domain.TaxDeduct, error) {
	deduct, err := s.taxRepository.GetTaxDeductByType(ctx, body.Type)
	if err != nil {
		return nil, err
	}

	if body.Amount.Cmp(deduct.MinAmount) < 0 {
		params := domain.NewFieldMessageList().Add(
			"amount", fmt.Sprintf("tax deduct type %s is less than %f", body.Type, deduct.MinAmount.InexactFloat64()),
		).Value()
		return nil, errorx.ErrAmountLessThanLimit.WithParams(params)
	}

	if body.Amount.Cmp(deduct.MaxAmount) > 0 {
		params := domain.NewFieldMessageList().Add(
			"amount", fmt.Sprintf("tax deduct type %s is more than %f", body.Type, deduct.MaxAmount.InexactFloat64()),
		).Value()
		return nil, errorx.ErrAmountMoreThanLimit.WithParams(params)
	}

	return s.taxRepository.UpdateTaxDeduct(ctx, body)
}
