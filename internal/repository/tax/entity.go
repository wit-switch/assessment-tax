package tax

import (
	"database/sql/driver"

	"github.com/wit-switch/assessment-tax/pkg/errorx"

	"github.com/shopspring/decimal"
)

type taxDeductType string

const (
	TaxDeductTypeDonation taxDeductType = "donation"
	TaxDeductTypePersonal taxDeductType = "personal"
	TaxDeductTypeKReceipt taxDeductType = "k-receipt"
)

func (e *taxDeductType) String() string {
	return string(*e)
}

func (e *taxDeductType) Scan(src any) error {
	switch s := src.(type) {
	case []byte:
		*e = taxDeductType(s)
	case string:
		*e = taxDeductType(s)
	default:
		return errorx.Errorf("unsupported scan type for TaxDeductType: %T", src)
	}
	return nil
}

type taxDeduct struct {
	Type      taxDeductType
	MinAmount decimal.Decimal
	MaxAmount decimal.Decimal
	Amount    decimal.Decimal
}

type nullTaxDeductType struct {
	TaxDeductType taxDeductType
	Valid         bool
}

func (ns *nullTaxDeductType) Scan(value any) error {
	if value == nil {
		ns.TaxDeductType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TaxDeductType.Scan(value)
}

func (ns nullTaxDeductType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.TaxDeductType.String(), nil
}

type getTaxDeduct struct {
	Type nullTaxDeductType
}
