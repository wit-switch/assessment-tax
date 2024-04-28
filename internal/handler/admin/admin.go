package admin

import "github.com/labstack/echo/v4"

// UpdatePersonalDeduct update tax deduct
// @Tags     admin
// @Accept   json
// @Produce  json
// @Security BasicAuth
// @Param    body body     updateTaxDeductRequest                true " "
// @Success  200  {object} updatePersonalDeductResponse          "Success"
// @Failure  400  {object} http.ResponseError[[]validator.Field] "Bad Request"
// @Failure  404  {object} http.ResponseError[string]            "Not Found"
// @Failure  500  {object} http.ResponseError[string]            "Internal Server Error"
// @Router   /admin/deductions/personal [post].
func (hdl *Handler) UpdatePersonalDeduct(
	c echo.Context,
	_ any,
	body updateTaxDeductRequest,
) (*updatePersonalDeductResponse, error) {
	req := hdl.dto.toUpdatePersonalDeductDomain(body)
	resp, err := hdl.taxService.UpdateTaxDeduct(c.Request().Context(), req)
	if err != nil {
		return nil, err
	}

	out := hdl.dto.toUpdatePersonalDeductResponse(*resp)
	return &out, nil
}
