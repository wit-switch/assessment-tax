package tax

import (
	"encoding/csv"
	"strings"

	"github.com/wit-switch/assessment-tax/pkg/errorx"

	"github.com/labstack/echo/v4"
)

// Calculate tax calculate
// @Tags    tax
// @Accept  json
// @Produce json
// @Param   body body     taxCalculateRequest                   true " "
// @Success 200  {object} taxCalculateResponse                  "Success"
// @Failure 400  {object} http.ResponseError[[]validator.Field] "Bad Request"
// @Failure 404  {object} http.ResponseError[string]            "Not Found"
// @Failure 500  {object} http.ResponseError[string]            "Internal Server Error"
// @Router  /tax/calculations [post].
func (hdl *Handler) Calculate(c echo.Context, _ any, req taxCalculateRequest) (*taxCalculateResponse, error) {
	body := hdl.dto.toTaxCalculateDomain(req)
	resp, err := hdl.taxService.Calculate(c.Request().Context(), body, true)
	if err != nil {
		return nil, err
	}

	out := hdl.dto.toTaxCalculateResponse(*resp)
	return &out, nil
}

// CalculateFromCSV tax calculate from csv file
// @Tags    tax
// @Accept  mpfd
// @Produce json
// @Param   taxFile formData file                                  true " "
// @Success 200     {object} texes                                 "Success"
// @Failure 400     {object} http.ResponseError[[]validator.Field] "Bad Request"
// @Failure 404     {object} http.ResponseError[string]            "Not Found"
// @Failure 500     {object} http.ResponseError[string]            "Internal Server Error"
// @Router  /tax/calculations/upload-csv [post].
func (hdl *Handler) CalculateFromCSV(
	c echo.Context,
	_, _ any,
) (*texes, error) {
	file, err := c.FormFile("taxFile")
	if err != nil {
		return nil, errorx.ErrValidationFail.WithError(err).WithParams([]map[string]any{
			{
				"field":   "taxFile",
				"message": "file in field taxFile is required",
			},
		})
	}

	if !strings.HasSuffix(file.Filename, ".csv") {
		return nil, errorx.ErrValidationFail.WithError(err).WithParams([]map[string]any{
			{
				"field":   "taxFile",
				"message": "file is not .csv",
			},
		})
	}

	formFile, err := file.Open()
	if err != nil {
		return nil, errorx.ErrValidationFail.WithError(err)
	}
	defer formFile.Close()

	r := csv.NewReader(formFile)
	resp, err := hdl.taxService.CalculateFromCSV(c.Request().Context(), *r)
	if err != nil {
		return nil, err
	}

	out := hdl.dto.toTaxes(resp)
	return &out, nil
}
