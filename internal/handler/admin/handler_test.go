package admin_test

import (
	"bytes"
	"context"
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/internal/handler/admin"
	httphdl "github.com/wit-switch/assessment-tax/internal/handler/http"
	"github.com/wit-switch/assessment-tax/mocks"
	"github.com/wit-switch/assessment-tax/pkg/errorx"
	"github.com/wit-switch/assessment-tax/pkg/validator"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Admin", func() {
	var (
		ctrl *gomock.Controller
		ctx  context.Context
		app  *echo.Echo
		hdl  *admin.Handler

		mockTaxService *mocks.MockTaxService

		taxRoute string
		// zero     decimal.Decimal
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockTaxService = mocks.NewMockTaxService(ctrl)
		hdl = admin.NewHandler(admin.Dependencies{
			TaxService: mockTaxService,
		})

		app = echo.New()
		validate, _ := validator.New()
		app.Validator = httphdl.NewValidator(validate)
		app.HTTPErrorHandler = httphdl.HTTPErrorHandler

		ctx = context.Background()

		taxRoute = "/admin/deductions"
	})

	AfterEach(func() {
		ctrl.Finish()
		ctrl = nil
	})

	Describe("update personal deduction", func() {
		var (
			route    string
			bodyJSON string

			mockTax *domain.TaxDeduct

			mockUpdateTaxDeduct *gomock.Call
		)

		BeforeEach(func() {
			// zero = decimal.NewFromFloat(0)
			route = taxRoute
			app.POST(route, httphdl.BindRoute(
				hdl.UpdatePersonalDeduct,
				httphdl.WithBodyParser(),
				httphdl.WithBodyValidator(),
			))

			bodyJSON = `{
				"amount": 70000.0
			}`

			mockTax = &domain.TaxDeduct{
				Type:      domain.TaxDeductTypePersonal,
				MinAmount: decimal.NewFromFloat(10000),
				MaxAmount: decimal.NewFromFloat(100000),
				Amount:    decimal.NewFromFloat(70000),
			}

			mockUpdateTaxDeduct = mockTaxService.EXPECT().
				UpdateTaxDeduct(ctx, domain.UpdateTaxDeduct{
					Type:   domain.TaxDeductTypePersonal,
					Amount: decimal.NewFromFloat(70000),
				}).MinTimes(0)
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

		})

		When("request is valid", func() {
			Context("with amount less than limit", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"amount": 700.0
					}`

					mockUpdateTaxDeduct = mockTaxService.EXPECT().
						UpdateTaxDeduct(ctx, domain.UpdateTaxDeduct{
							Type:   domain.TaxDeductTypePersonal,
							Amount: decimal.NewFromFloat(700),
						})
				})
				It("should return 400 bad request", func() {
					mockUpdateTaxDeduct.Return(nil, errorx.ErrAmountLessThanLimit)

					code, respBody := request(
						route,
						bytes.NewBufferString(bodyJSON),
						app,
					)
					Expect(http.StatusBadRequest).To(Equal(code))

					actual, err := compacJSON(respBody)
					Expect(err).NotTo(HaveOccurred())
					expectedResp := `{
						"code":"400102",
						"message":"amount less than limit"
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			Context("with amount more than limit", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"amount": 70000000.0
					}`

					mockUpdateTaxDeduct = mockTaxService.EXPECT().
						UpdateTaxDeduct(ctx, domain.UpdateTaxDeduct{
							Type:   domain.TaxDeductTypePersonal,
							Amount: decimal.NewFromFloat(70000000),
						})
				})
				It("should return 400 bad request", func() {
					mockUpdateTaxDeduct.Return(nil, errorx.ErrAmountMoreThanLimit)

					code, respBody := request(
						route,
						bytes.NewBufferString(bodyJSON),
						app,
					)
					Expect(http.StatusBadRequest).To(Equal(code))

					actual, err := compacJSON(respBody)
					Expect(err).NotTo(HaveOccurred())
					expectedResp := `{
						"code":"400101",
						"message":"amount more than limit"
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})

				It("should return 200 ok", func() {
					mockUpdateTaxDeduct.Return(mockTax, nil)

					code, respBody := request(
						route,
						bytes.NewBufferString(bodyJSON),
						app,
					)
					Expect(http.StatusOK).To(Equal(code))

					actual, err := compacJSON(respBody)
					Expect(err).NotTo(HaveOccurred())
					expectedResp := `{
						"personalDeduction": 70000
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})
		})
	})

	Describe("update k-receipt deduction", func() {
		var (
			route    string
			bodyJSON string

			mockTax *domain.TaxDeduct

			mockUpdateTaxDeduct *gomock.Call
		)

		BeforeEach(func() {
			// zero = decimal.NewFromFloat(0)
			route = taxRoute
			app.POST(route, httphdl.BindRoute(
				hdl.UpdateKReceiptDeduct,
				httphdl.WithBodyParser(),
				httphdl.WithBodyValidator(),
			))

			bodyJSON = `{
				"amount": 70000.0
			}`

			mockTax = &domain.TaxDeduct{
				Type:      domain.TaxDeductTypeKReceipt,
				MinAmount: decimal.NewFromFloat(1),
				MaxAmount: decimal.NewFromFloat(100000),
				Amount:    decimal.NewFromFloat(70000),
			}

			mockUpdateTaxDeduct = mockTaxService.EXPECT().
				UpdateTaxDeduct(ctx, domain.UpdateTaxDeduct{
					Type:   domain.TaxDeductTypeKReceipt,
					Amount: decimal.NewFromFloat(70000),
				}).MinTimes(0)
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

		})

		When("request is valid", func() {
			Context("with amount less than limit", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"amount": 0.0
					}`

					mockUpdateTaxDeduct = mockTaxService.EXPECT().
						UpdateTaxDeduct(ctx, domain.UpdateTaxDeduct{
							Type:   domain.TaxDeductTypeKReceipt,
							Amount: decimal.NewFromFloat(0),
						})
				})
				It("should return 400 bad request", func() {
					mockUpdateTaxDeduct.Return(nil, errorx.ErrAmountLessThanLimit)

					code, respBody := request(
						route,
						bytes.NewBufferString(bodyJSON),
						app,
					)
					Expect(http.StatusBadRequest).To(Equal(code))

					actual, err := compacJSON(respBody)
					Expect(err).NotTo(HaveOccurred())
					expectedResp := `{
						"code":"400102",
						"message":"amount less than limit"
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			Context("with amount more than limit", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"amount": 70000000.0
					}`

					mockUpdateTaxDeduct = mockTaxService.EXPECT().
						UpdateTaxDeduct(ctx, domain.UpdateTaxDeduct{
							Type:   domain.TaxDeductTypeKReceipt,
							Amount: decimal.NewFromFloat(70000000),
						})
				})
				It("should return 400 bad request", func() {
					mockUpdateTaxDeduct.Return(nil, errorx.ErrAmountMoreThanLimit)

					code, respBody := request(
						route,
						bytes.NewBufferString(bodyJSON),
						app,
					)
					Expect(http.StatusBadRequest).To(Equal(code))

					actual, err := compacJSON(respBody)
					Expect(err).NotTo(HaveOccurred())
					expectedResp := `{
						"code":"400101",
						"message":"amount more than limit"
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})

				It("should return 200 ok", func() {
					mockUpdateTaxDeduct.Return(mockTax, nil)

					code, respBody := request(
						route,
						bytes.NewBufferString(bodyJSON),
						app,
					)
					Expect(http.StatusOK).To(Equal(code))

					actual, err := compacJSON(respBody)
					Expect(err).NotTo(HaveOccurred())
					expectedResp := `{
						"kReceipt": 70000
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})
		})
	})
})
