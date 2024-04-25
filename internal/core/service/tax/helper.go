package tax

import (
	"fmt"

	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/pkg/errorx"
)

// GetTaxDeductByType get and validate tax deduct type should have in db.
func GetTaxDeductByType(
	taxDeductMap map[domain.TaxDeductType]domain.TaxDeduct,
	taxDeductType domain.TaxDeductType,
) (*domain.TaxDeduct, error) {
	out, ok := taxDeductMap[taxDeductType]
	if !ok {
		params := domain.NewFieldMessageList().Add(
			"type", fmt.Sprintf("tax deduct type %s is not found", taxDeductType),
		).Value()
		return nil, errorx.ErrTaxDeductNotFound.WithParams(params)
	}

	return &out, nil
}
