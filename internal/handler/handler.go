package handler

import (
	"github.com/wit-switch/assessment-tax/internal/core/service"
	"github.com/wit-switch/assessment-tax/internal/handler/admin"
	"github.com/wit-switch/assessment-tax/internal/handler/tax"
)

type Dependencies struct {
	Services *service.Services
}

type Handler struct {
	Admin *admin.Handler
	Tax   *tax.Handler
}

func New(deps Dependencies) *Handler {
	return &Handler{
		Admin: admin.NewHandler(admin.Dependencies{
			TaxService: deps.Services.Tax,
		}),
		Tax: tax.NewHandler(tax.Dependencies{
			TaxService: deps.Services.Tax,
		}),
	}
}
