package tax

import (
	"github.com/wit-switch/assessment-tax/internal/core/domain"

	"github.com/shopspring/decimal"
)

type dto struct{}

func (d *dto) toTaxCalculateDomain(a taxCalculateRequest) domain.TaxCalculate {
	return domain.TaxCalculate{
		TotalIncome: decimal.NewFromFloat(a.TotalIncome),
		Wht:         decimal.NewFromFloat(a.Wht),
		Allowances:  d.toAllowancesDomain(a.Allowances),
	}
}

func (d *dto) toAllowanceDomain(a allowance) domain.Allowance {
	return domain.Allowance{
		AllowanceType: domain.TaxDeductType(a.AllowanceType),
		Amount:        decimal.NewFromFloat(a.Amount),
	}
}

func (d *dto) toAllowancesDomain(a []allowance) []domain.Allowance {
	out := make([]domain.Allowance, len(a))
	for i, v := range a {
		out[i] = d.toAllowanceDomain(v)
	}

	return out
}

func (d *dto) toTaxCalculateResponse(a domain.Tax) taxCalculateResponse {
	return taxCalculateResponse{
		Tax: a.Tax.InexactFloat64(),
	}
}