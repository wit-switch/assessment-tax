package admin

type updateTaxDeductRequest struct {
	Amount float64 `json:"amount" validate:"gte=0"`
}

type updatePersonalDeductResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}
