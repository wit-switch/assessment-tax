package tax

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wit-switch/assessment-tax/internal/core/port"
)

type Dependencies struct {
	DB *pgxpool.Pool
}

type Repositories struct {
	db  *pgxpool.Pool
	dto dto
}

var _ port.TaxRepository = (*Repositories)(nil)

func NewRepository(deps Dependencies) *Repositories {
	return &Repositories{
		db:  deps.DB,
		dto: dto{},
	}
}
