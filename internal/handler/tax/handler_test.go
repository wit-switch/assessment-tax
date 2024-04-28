package tax_test

import (
	"bytes"
	"context"
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/wit-switch/assessment-tax/internal/core/domain"
	httphdl "github.com/wit-switch/assessment-tax/internal/handler/http"
	"github.com/wit-switch/assessment-tax/internal/handler/tax"
	"github.com/wit-switch/assessment-tax/mocks"
	"github.com/wit-switch/assessment-tax/pkg/errorx"
	"github.com/wit-switch/assessment-tax/pkg/validator"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Tax", func() {
	var (
		ctrl *gomock.Controller
		ctx  context.Context
		app  *echo.Echo
		hdl  *tax.Handler

		mockTaxService *mocks.MockTaxService

		taxRoute string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockTaxService = mocks.NewMockTaxService(ctrl)
		hdl = tax.NewHandler(tax.Dependencies{
			TaxService: mockTaxService,
		})

		app = echo.New()
		app.Validator = httphdl.NewValidator(validator.New())
		app.HTTPErrorHandler = httphdl.HTTPErrorHandler

		ctx = context.Background()

		taxRoute = "/tax/calculations"
	})

	AfterEach(func() {
		ctrl.Finish()
		ctrl = nil
	})

	Describe("calculate tax", func() {
		var (
			route    string
			bodyJSON string

			mockTax *domain.Tax

			mockCalculate *gomock.Call
		)

		BeforeEach(func() {
			route = taxRoute
			app.POST(route, httphdl.BindRoute(
				hdl.Calculate,
				httphdl.WithBodyParser(),
				httphdl.WithBodyValidator(),
			))

			bodyJSON = `{
				"totalIncome": 500000.0,
				"wht": 0.0,
				"allowances": [
					{
						"allowanceType": "donation",
						"amount": 0.0
					}
				]
			}`

			mockTax = &domain.Tax{
				Tax: decimal.NewFromFloat(29000),
			}

			mockCalculate = mockTaxService.EXPECT().
				Calculate(ctx, domain.TaxCalculate{
					TotalIncome: decimal.NewFromFloat(500000),
					Wht:         decimal.NewFromFloat(0),
					Allowances: []domain.Allowance{
						{
							AllowanceType: domain.TaxDeductTypeDonation,
							Amount:        decimal.NewFromFloat(0),
						},
					},
				}).
				MinTimes(0)
		})

		When("request is not valid", func() {
			Context("with body is invalid", func() {
				BeforeEach(func() {
					bodyJSON = `x`
				})
				It("should return 400 validate error", func() {
					code, respBody := request(
						route,
						bytes.NewBufferString(bodyJSON),
						app,
					)
					Expect(http.StatusBadRequest).To(Equal(code))

					actual, err := compacJSON(respBody)
					Expect(err).NotTo(HaveOccurred())

					expectedResp := `{
						"code":"400001",
						"message":"validation error"
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			Context("with totalIncome less than zero", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"totalIncome": -500000.0,
						"wht": 0.0,
						"allowances": [
							{
								"allowanceType": "donation",
								"amount": 0.0
							}
						]
					}`
				})
				It("should return 400 validate error with field error", func() {
					code, respBody := request(
						route,
						bytes.NewBufferString(bodyJSON),
						app,
					)
					Expect(http.StatusBadRequest).To(Equal(code))

					actual, err := compacJSON(respBody)
					Expect(err).NotTo(HaveOccurred())

					expectedResp := `{
						"code":"400001",
						"message":"validation error",
						"errors":[
							{
								"field":"totalIncome",
								"message":"totalIncome should greater than or equal 0"
							}
						]
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})
		})

		When("request is valid", func() {
			Context("with failed to calculate tax", func() {
				errUnknown := errorx.New("error")
				It("should return 500 unknown error", func() {
					mockCalculate.Return(nil, errUnknown)

					code, respBody := request(
						route,
						bytes.NewBufferString(bodyJSON),
						app,
					)
					Expect(http.StatusInternalServerError).To(Equal(code))

					actual, err := compacJSON(respBody)
					Expect(err).NotTo(HaveOccurred())

					expectedResp := `{
						"code":"500000",
						"message":"error"
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			It("should return tax", func() {
				mockCalculate.Return(mockTax, nil)

				code, respBody := request(
					route,
					bytes.NewBufferString(bodyJSON),
					app,
				)
				Expect(http.StatusOK).To(Equal(code))

				actual, err := compacJSON(respBody)
				Expect(err).NotTo(HaveOccurred())

				expectedResp := `{
					"tax":29000
				}`
				expected, err := compacJSON(expectedResp)
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
