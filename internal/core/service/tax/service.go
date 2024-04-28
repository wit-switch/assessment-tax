package tax

import (
	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/internal/core/port"

	"github.com/shopspring/decimal"
)

type Dependencies struct {
	TaxRepository port.TaxRepository
}

type Services struct {
	taxRepository port.TaxRepository
}

var _ port.TaxService = (*Services)(nil)

func NewService(deps Dependencies) *Services {
	return &Services{
		taxRepository: deps.TaxRepository,
	}
}

var taxRates = []domain.TaxRate{
	{
		Description: "2,000,001 ขึ้นไป",
		Rate:        decimal.NewFromFloat(0.35),
		MinIncome:   decimal.NewFromInt(2000000),
	},
	{
		Description: "1,000,001-2,000,000",
		Rate:        decimal.NewFromFloat(0.2),
		MinIncome:   decimal.NewFromInt(1000000),
	},
	{
		Description: "500,001-1,000,000",
		Rate:        decimal.NewFromFloat(0.15),
		MinIncome:   decimal.NewFromInt(500000),
	},
	{
		Description: "150,001-500,000",
		Rate:        decimal.NewFromFloat(0.1),
		MinIncome:   decimal.NewFromInt(150000),
	},
	{
		Description: "0-150,000",
		Rate:        decimal.NewFromFloat(0),
		MinIncome:   decimal.NewFromInt(0),
	},
}
