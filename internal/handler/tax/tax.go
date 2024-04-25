package tax

import "github.com/labstack/echo/v4"

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
	resp, err := hdl.taxService.Calculate(c.Request().Context(), body)
	if err != nil {
		return nil, err
	}

	out := hdl.dto.toTaxCalculateResponse(*resp)
	return &out, nil
}
