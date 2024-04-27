package domain

import "github.com/shopspring/decimal"

type Decimal struct {
	decimal.Decimal
}

type CSVFile struct {
	TotalIncome Decimal `csv:"totalIncome"`
	Wht         Decimal `csv:"wht"`
	Donation    Decimal `csv:"donation"`
}

func (c *Decimal) UnmarshalCSV(csv string) error {
	value, err := decimal.NewFromString(csv)
	if err != nil {
		return err
	}

	c.Decimal = value
	return nil
}
