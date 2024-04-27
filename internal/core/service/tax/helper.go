package tax

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/shopspring/decimal"
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

func GetAllowanceAmount(
	ctx context.Context,
	taxDeduct *domain.TaxDeduct,
	allowanceAmount decimal.Decimal,
) (decimal.Decimal, error) {
	amount := allowanceAmount
	if amount.Cmp(taxDeduct.Amount) > 0 {
		slog.WarnContext(ctx, "amount more than limit",
			slog.String("type", domain.TaxDeductTypeDonation.String()),
			slog.Float64("amount", amount.InexactFloat64()),
			slog.Float64("max_amount", taxDeduct.MaxAmount.InexactFloat64()),
		)

		amount = taxDeduct.Amount
	}

	if amount.Cmp(taxDeduct.MinAmount) < 0 {
		params := domain.NewFieldMessageList().Add(
			"allowanceType", fmt.Sprintf("allowance type %s is less than limit", taxDeduct.Type),
		).Value()
		return amount, errorx.ErrAmountLessThanLimit.WithParams(params)
	}

	return amount, nil
}
