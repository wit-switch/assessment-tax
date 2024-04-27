package tax

import (
	"context"
	"time"

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

const getTaxDeductByTypeSQL = `-- name: GetTaxDeductByType :one
SELECT type, min_amount, max_amount, amount FROM tax_deduct
WHERE 
type = $1
LIMIT 1
`

func (repo *Repositories) GetTaxDeductByType(ctx context.Context, q domain.TaxDeductType) (*domain.TaxDeduct, error) {
	row := repo.db.QueryRow(ctx, getTaxDeductByTypeSQL, taxDeductType(q))
	var i taxDeduct
	err := row.Scan(
		&i.Type,
		&i.MinAmount,
		&i.MaxAmount,
		&i.Amount,
	)
	if err != nil {
		return nil, errorInterceptor(err)
	}

	out := repo.dto.toTaxDeductDomain(i)
	return &out, nil
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

const updateTaxDeductSQL = `-- name: UpdateTaxDeduct :one
UPDATE tax_deduct 
SET 
amount = $2,
updated_at = $3
WHERE type = $1
RETURNING type, min_amount, max_amount, amount
`

func (repo *Repositories) UpdateTaxDeduct(ctx context.Context, u domain.UpdateTaxDeduct) (*domain.TaxDeduct, error) {
	args := repo.dto.tpUpdateTaxDeduct(u)
	if args.UpdatedAt.IsZero() {
		args.UpdatedAt = time.Now().UTC()
	}

	row := repo.db.QueryRow(ctx, updateTaxDeductSQL, args.Type, args.Amount, args.UpdatedAt)
	var i taxDeduct
	err := row.Scan(
		&i.Type,
		&i.MinAmount,
		&i.MaxAmount,
		&i.Amount,
	)
	if err != nil {
		return nil, errorInterceptor(err)
	}

	out := repo.dto.toTaxDeductDomain(i)
	return &out, nil
}
