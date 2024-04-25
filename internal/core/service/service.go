package service

import (
	"github.com/wit-switch/assessment-tax/internal/core/port"
	"github.com/wit-switch/assessment-tax/internal/core/service/tax"
	"github.com/wit-switch/assessment-tax/internal/repository"
)

type Dependencies struct {
	Repositories *repository.Repositories
}

type Services struct {
	Tax port.TaxService
}

func New(deps Dependencies) *Services {
	return &Services{
		Tax: tax.NewService(tax.Dependencies{
			TaxRepository: deps.Repositories.Tax,
		}),
	}
}
