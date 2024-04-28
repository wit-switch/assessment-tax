package tax_test

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

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
		zero     decimal.Decimal
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockTaxService = mocks.NewMockTaxService(ctrl)
		hdl = tax.NewHandler(tax.Dependencies{
			TaxService: mockTaxService,
		})

		app = echo.New()
		validate, _ := validator.New()
		app.Validator = httphdl.NewValidator(validate)
		app.HTTPErrorHandler = httphdl.HTTPErrorHandler

		ctx = context.Background()

		taxRoute = "/tax/calculations"
		zero = decimal.NewFromFloat(0)
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
				TaxLevel: []domain.TaxLevel{
					{
						Level: "0-150,000",
						Tax:   zero,
					},
					{
						Level: "150,001-500,000",
						Tax:   decimal.NewFromFloat(29000),
					},
					{
						Level: "500,001-1,000,000",
						Tax:   zero,
					},
					{
						Level: "1,000,001-2,000,000",
						Tax:   zero,
					},
					{
						Level: "2,000,001 ขึ้นไป",
						Tax:   zero,
					},
				},
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
				}, true).
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
							},
							{
								"field":"wht",
								"message":"wht should less than or equal totalIncome"
							}
						]
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			Context("with wht less than zero", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"totalIncome": 1.0,
						"wht": -2.0,
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
								"field":"wht",
								"message":"wht should greater than or equal 0"
							}
						]
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			Context("with wht less than or equal totalIncome", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"totalIncome": 1.0,
						"wht": 2.0,
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
								"field":"wht",
								"message":"wht should less than or equal totalIncome"
							}
						]
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			Context("with allowances is empty", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"totalIncome": 500000.0,
						"wht": 0.0
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
								"field":"allowances",
								"message":"allowances is required"
							}
						]
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			Context("with allowances is not unique", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"totalIncome": 500000.0,
						"wht": 0.0,
						"allowances": [
							{
								"allowanceType": "donation",
								"amount": 0.0
							},
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
								"field":"allowances",
								"message":"allowances is unique"
							}
						]
					}`
					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			Context("with allowanceType is not valide", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"totalIncome": 500000.0,
						"wht": 0.0,
						"allowances": [
							{
								"allowanceType": "xxx",
								"amount": 0.0
							},
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
								"field":"allowanceType",
								"message":"allowanceType [xxx] is not valide"
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

			Context("with have tax refund", func() {
				BeforeEach(func() {
					bodyJSON = `{
						"totalIncome": 450000.0,
						"wht": 25000.0,
						"allowances": [
							{
								"allowanceType": "donation",
								"amount": 0.0
							}
						]
					}`

					mockCalculate = mockTaxService.EXPECT().
						Calculate(ctx, domain.TaxCalculate{
							TotalIncome: decimal.NewFromFloat(450000),
							Wht:         decimal.NewFromFloat(25000),
							Allowances: []domain.Allowance{
								{
									AllowanceType: domain.TaxDeductTypeDonation,
									Amount:        decimal.NewFromFloat(0),
								},
							},
						}, true)

					mockTax.Tax = zero
					mockTax.TaxRefund = decimal.NewFromFloat(1000)
					mockTax.TaxLevel = []domain.TaxLevel{
						{
							Level: "0-150,000",
							Tax:   zero,
						},
						{
							Level: "150,001-500,000",
							Tax:   decimal.NewFromFloat(24000),
						},
						{
							Level: "500,001-1,000,000",
							Tax:   zero,
						},
						{
							Level: "1,000,001-2,000,000",
							Tax:   zero,
						},
						{
							Level: "2,000,001 ขึ้นไป",
							Tax:   zero,
						},
					}
				})
				It("should return tax refund", func() {
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
						"tax":0,
						"taxRefund":1000,
						"taxLevel": [
							{
								"level": "0-150,000",
								"tax": 0
							},
							{
								"level": "150,001-500,000",
								"tax": 24000
							},
							{
								"level": "500,001-1,000,000",
								"tax": 0
							},
							{
								"level": "1,000,001-2,000,000",
								"tax": 0
							},
							{
								"level": "2,000,001 ขึ้นไป",
								"tax": 0
							}
						]
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
					"tax":29000,
					"taxLevel": [
					{
						"level": "0-150,000",
						"tax": 0
					},
					{
						"level": "150,001-500,000",
						"tax": 29000
					},
					{
						"level": "500,001-1,000,000",
						"tax": 0
					},
					{
						"level": "1,000,001-2,000,000",
						"tax": 0
					},
					{
						"level": "2,000,001 ขึ้นไป",
						"tax": 0
					}
				]
				}`
				expected, err := compacJSON(expectedResp)
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(Equal(expected))
			})
		})
	})

	Describe("calculate tax from csv", func() {
		var (
			route string
		)

		BeforeEach(func() {
			route = fmt.Sprintf("%s/upload-csv", taxRoute)
			app.POST(route, httphdl.BindRoute(
				hdl.CalculateFromCSV,
			))
		})

		When("request is valid", func() {
			Context("with wrong file name", func() {
				It("should return  an error", func() {
					buffer := new(bytes.Buffer)
					writer := multipart.NewWriter(buffer)
					formFile, _ := writer.CreateFormFile("File", "taxes.csv")

					csvData := `totalIncome,wht,donation
500000,0,0
600000,40000,20000
750000,50000,15000
`
					_, _ = formFile.Write([]byte(csvData))
					_ = writer.Close()

					mockTaxService.EXPECT().
						CalculateFromCSV(ctx, gomock.Any()).
						Times(0)

					req := httptest.NewRequest(http.MethodPost, route, buffer)
					req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
					rec := httptest.NewRecorder()
					// custom error handler only plays with real request
					app.ServeHTTP(rec, req)

					Expect(http.StatusBadRequest).To(Equal(rec.Code))

					actual, err := compacJSON(rec.Body.String())
					Expect(err).NotTo(HaveOccurred())

					expectedResp := `{
						"code":"400001",
						"message":"validation error",
						"errors":[
							{
								"field":"taxFile",
								"message":"file in field taxFile is required"
							}
						]
					}`

					expected, err := compacJSON(expectedResp)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(expected))
				})
			})

			It("should return taxes", func() {
				buffer := new(bytes.Buffer)
				writer := multipart.NewWriter(buffer)
				formFile, _ := writer.CreateFormFile("taxFile", "taxes.csv")

				csvData := `totalIncome,wht,donation
500000,0,0
600000,40000,20000
750000,50000,15000
`
				_, _ = formFile.Write([]byte(csvData))
				_ = writer.Close()

				mockTaxService.EXPECT().
					CalculateFromCSV(ctx, gomock.Any()).
					Times(1).
					Return([]domain.TaxCSV{
						{
							TotalIncome: decimal.NewFromFloat(500000),
							Tax:         decimal.NewFromFloat(28000),
							TaxRefund:   zero,
						},
						{
							TotalIncome: decimal.NewFromFloat(600000),
							Tax:         zero,
							TaxRefund:   decimal.NewFromFloat(3500),
						},
						{
							TotalIncome: decimal.NewFromFloat(750000),
							Tax:         decimal.NewFromFloat(9750),
							TaxRefund:   zero,
						},
					}, nil)

				req := httptest.NewRequest(http.MethodPost, route, buffer)
				req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
				rec := httptest.NewRecorder()
				// custom error handler only plays with real request
				app.ServeHTTP(rec, req)

				Expect(http.StatusOK).To(Equal(rec.Code))

				actual, err := compacJSON(rec.Body.String())
				Expect(err).NotTo(HaveOccurred())

				expectedResp := `{
					"texes": [
						{
							"totalIncome": 500000,
							"tax": 28000
						},
						{
							"totalIncome": 600000,
							"tax": 0,
							"taxRefund": 3500
						},
						{
							"totalIncome": 750000,
							"tax": 9750
						}
					]
				}`

				expected, err := compacJSON(expectedResp)
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
