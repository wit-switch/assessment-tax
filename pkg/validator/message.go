package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

func getValidatorMessage(v validator.FieldError) string {
	switch v.Tag() {
	case "gt":
		return fmt.Sprintf("%s should greater than %s", v.Field(), v.Param())
	case "gtcsfield":
		return fmt.Sprintf("%s should greater than %s", v.Field(), strcase.ToLowerCamel(v.Param()))
	case "gte":
		return fmt.Sprintf("%s should greater than or equal %s", v.Field(), v.Param())
	case "gtecsfield":
		return fmt.Sprintf("%s should greater than or equal %s", v.Field(), strcase.ToLowerCamel(v.Param()))
	case "lt":
		return fmt.Sprintf("%s should less than %s", v.Field(), v.Param())
	case "ltcsfield":
		return fmt.Sprintf("%s should less than %s", v.Field(), strcase.ToLowerCamel(v.Param()))
	case "lte":
		return fmt.Sprintf("%s should less than or equal %s", v.Field(), v.Param())
	case "ltecsfield":
		return fmt.Sprintf("%s should less than or equal %s", v.Field(), strcase.ToLowerCamel(v.Param()))
	default:
		return fmt.Sprintf("%s is %s", v.Field(), v.Tag())
	}
}
