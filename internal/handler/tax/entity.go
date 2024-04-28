package tax

type taxCalculateRequest struct {
	TotalIncome float64     `json:"totalIncome" validate:"gte=0" example:"500000.0"`
	Wht         float64     `json:"wht"`
	Allowances  []allowance `json:"allowances"`
}

type allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type taxCalculateResponse struct {
	Tax float64 `json:"tax"`
}
