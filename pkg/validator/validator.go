package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	customeAllowanceType = "with-allowance-type"

	allowanceTypeDonation = "donation"
	allowanceTypeKReceipt = "k-receipt"
)

type Validator interface {
	Struct(s any) error
}

type wrapValidator struct {
	*validator.Validate
}

func New() (Validator, error) {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	err := v.RegisterValidation(customeAllowanceType, validateAllowanceType)
	if err != nil {
		return nil, err
	}

	return &wrapValidator{v}, nil
}

func (w *wrapValidator) Struct(s any) error {
	if err := w.Validate.Struct(s); err != nil {
		return &FieldError{err: err}
	}

	return nil
}

func validateAllowanceType(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	return v == allowanceTypeDonation ||
		v == allowanceTypeKReceipt
}
