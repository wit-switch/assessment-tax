package domain

import (
	"github.com/shopspring/decimal"
)

type TaxDeductType string

const (
	TaxDeductTypeDonation TaxDeductType = "donation"
	TaxDeductTypePersonal TaxDeductType = "personal"
	TaxDeductTypeKReceipt TaxDeductType = "k-receipt"
)

func (d TaxDeductType) String() string {
	return string(d)
}

type TaxRate struct {
	Description string
	Rate        decimal.Decimal
	MinIncome   decimal.Decimal
}

type TaxCalculate struct {
	TotalIncome decimal.Decimal
	Wht         decimal.Decimal
	Allowances  []Allowance
}

type Allowance struct {
	AllowanceType TaxDeductType
	Amount        decimal.Decimal
}

type Tax struct {
	Tax       decimal.Decimal
	TaxRefund decimal.Decimal
	TaxLevel  []TaxLevel
}

type TaxLevel struct {
	Level string
	Tax   decimal.Decimal
}

type TaxDeduct struct {
	Type      TaxDeductType
	MinAmount decimal.Decimal
	MaxAmount decimal.Decimal
	Amount    decimal.Decimal
}

type GetTaxDeduct struct {
	Type TaxDeductType
}
