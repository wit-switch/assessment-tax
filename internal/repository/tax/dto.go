package tax

import (
	"github.com/wit-switch/assessment-tax/internal/core/domain"
)

type dto struct{}

func (d *dto) toTaxDeductDomain(a taxDeduct) domain.TaxDeduct {
	return domain.TaxDeduct{
		Type:      domain.TaxDeductType(a.Type),
		MinAmount: a.MinAmount,
		MaxAmount: a.MaxAmount,
		Amount:    a.Amount,
	}
}

func (d *dto) toTaxDeductsDomain(a []taxDeduct) []domain.TaxDeduct {
	out := make([]domain.TaxDeduct, len(a))
	for i, v := range a {
		out[i] = d.toTaxDeductDomain(v)
	}

	return out
}

func (d *dto) toGetTaxDeduct(a domain.GetTaxDeduct) getTaxDeduct {
	var deductType nullTaxDeductType

	if a.Type != "" {
		deductType.TaxDeductType = taxDeductType(a.Type)
		deductType.Valid = true
	}

	return getTaxDeduct{
		Type: deductType,
	}
}

func (d *dto) tpUpdateTaxDeduct(a domain.UpdateTaxDeduct) updateTaxDeduct {
	return updateTaxDeduct{
		Type:   taxDeductType(a.Type),
		Amount: a.Amount,
	}
}
