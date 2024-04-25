package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wit-switch/assessment-tax/internal/core/port"
	"github.com/wit-switch/assessment-tax/internal/repository/tax"
)

type Dependencies struct {
	DB *pgxpool.Pool
}

type Repositories struct {
	Tax port.TaxRepository
}

func New(deps Dependencies) *Repositories {
	return &Repositories{
		Tax: tax.NewRepository(tax.Dependencies{
			DB: deps.DB,
		}),
	}
}
