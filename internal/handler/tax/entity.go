package tax

type taxCalculateRequest struct {
	TotalIncome float64     `json:"totalIncome" validate:"gte=0" example:"500000.0"`
	Wht         float64     `json:"wht" validate:"gte=0,ltecsfield=TotalIncome" example:"25000.0"`
	Allowances  []allowance `json:"allowances" validate:"required,gt=0,unique=AllowanceType,dive"`
}

type allowance struct {
	AllowanceType string  `json:"allowanceType" validate:"required,with-allowance-type" example:"donation,k-receipt"`
	Amount        float64 `json:"amount" validate:"gte=0" example:"200000.0,100000.0"`
}

type taxCalculateResponse struct {
	Tax       float64    `json:"tax"`
	TaxRefund float64    `json:"taxRefund,omitempty"`
	TaxLevel  []taxLevel `json:"taxLevel"`
}

type taxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type texes struct {
	Texes []taxCSV `json:"texes"`
}

type taxCSV struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
	TaxRefund   float64 `json:"taxRefund,omitempty"`
}
