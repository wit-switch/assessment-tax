package handler

import (
	"github.com/wit-switch/assessment-tax/internal/core/service"
	"github.com/wit-switch/assessment-tax/internal/handler/tax"
)

type Dependencies struct {
	Services *service.Services
}

type Handler struct {
	Tax *tax.Handler
}

func New(deps Dependencies) *Handler {
	return &Handler{
		Tax: tax.NewHandler(tax.Dependencies{
			TaxService: deps.Services.Tax,
		}),
	}
}
