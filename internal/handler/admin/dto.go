package admin

import (
	"github.com/wit-switch/assessment-tax/internal/core/domain"

	"github.com/shopspring/decimal"
)

type dto struct{}

func (d *dto) toUpdatePersonalDeductDomain(a updateTaxDeductRequest) domain.UpdateTaxDeduct {
	return domain.UpdateTaxDeduct{
		Type:   domain.TaxDeductTypePersonal,
		Amount: decimal.NewFromFloat(a.Amount),
	}
}

func (d *dto) toUpdatePersonalDeductResponse(a domain.TaxDeduct) updatePersonalDeductResponse {
	return updatePersonalDeductResponse{
		PersonalDeduction: a.Amount.InexactFloat64(),
	}
}
