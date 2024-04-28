package admin

import "github.com/labstack/echo/v4"

// UpdatePersonalDeduct update personal deduction
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

// UpdateKReceiptDeduct update k-receipt deduction
// @Tags     admin
// @Accept   json
// @Produce  json
// @Security BasicAuth
// @Param    body body     updateTaxDeductRequest                true " "
// @Success  200  {object} updateKReceiptDeductResponse          "Success"
// @Failure  400  {object} http.ResponseError[[]validator.Field] "Bad Request"
// @Failure  404  {object} http.ResponseError[string]            "Not Found"
// @Failure  500  {object} http.ResponseError[string]            "Internal Server Error"
// @Router   /admin/deductions/k-receipt [post].
func (hdl *Handler) UpdateKReceiptDeduct(
	c echo.Context,
	_ any,
	body updateTaxDeductRequest,
) (*updateKReceiptDeductResponse, error) {
	resp, err := hdl.taxService.UpdateTaxDeduct(c.Request().Context(), hdl.dto.toUpdateKReceiptDeductDomain(body))
	if err != nil {
		return nil, err
	}

	out := hdl.dto.toUpdateKReceiptDeductResponse(*resp)
	return &out, nil
}
