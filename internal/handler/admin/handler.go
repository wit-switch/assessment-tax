package admin

import (
	"github.com/wit-switch/assessment-tax/internal/core/port"
)

type Dependencies struct {
	TaxService port.TaxService
}

type Handler struct {
	taxService port.TaxService
	dto        dto
}

func NewHandler(deps Dependencies) *Handler {
	return &Handler{
		taxService: deps.TaxService,
		dto:        dto{},
	}
}
