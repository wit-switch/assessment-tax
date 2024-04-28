package errorx

type ErrCode int

func (e ErrCode) Int() int {
	return int(e)
}

var (
	CodeUnknown        ErrCode
	CodeValidationFail ErrCode = 1
	CodeUnauthorized   ErrCode = 2

	CodeTaxDeductNotFound ErrCode = 100
)

var (
	ErrValidationFail = NewInternalErr[any](CodeValidationFail)
	ErrUnauthorized   = NewInternalErr[any](CodeUnauthorized)

	ErrTaxDeductNotFound = NewInternalErr[any](CodeTaxDeductNotFound)
)
