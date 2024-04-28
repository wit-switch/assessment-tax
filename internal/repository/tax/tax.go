package tax

import (
	"context"

	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/pkg/errorx"

	"github.com/jackc/pgx/v5"
)

func errorInterceptor(err error) error {
	if errorx.Is(err, pgx.ErrNoRows) {
		return errorx.ErrTaxDeductNotFound.WithError(err)
	}

	return err
}

const listTaxDeductSQL = `-- name: ListTaxDeduct :many
SELECT type, min_amount, max_amount, amount FROM tax_deduct
WHERE 
type = coalesce($1, type)
ORDER BY type
`

func (repo *Repositories) ListTaxDeduct(ctx context.Context, q domain.GetTaxDeduct) ([]domain.TaxDeduct, error) {
	query := repo.dto.toGetTaxDeduct(q)
	rows, err := repo.db.Query(ctx, listTaxDeductSQL, query.Type)
	if err != nil {
		return nil, errorInterceptor(err)
	}
	defer rows.Close()

	var items []taxDeduct
	for rows.Next() {
		var i taxDeduct
		if errScan := rows.Scan(
			&i.Type,
			&i.MinAmount,
			&i.MaxAmount,
			&i.Amount,
		); errScan != nil {
			return nil, errorInterceptor(errScan)
		}
		items = append(items, i)
	}
	if rErr := rows.Err(); rErr != nil {
		return nil, errorInterceptor(rErr)
	}

	if len(items) == 0 {
		return nil, errorx.ErrTaxDeductNotFound
	}

	out := repo.dto.toTaxDeductsDomain(items)
	return out, nil
}
