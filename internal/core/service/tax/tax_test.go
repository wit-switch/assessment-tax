package tax_test

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/internal/core/service/tax"
	"github.com/wit-switch/assessment-tax/mocks"
	"github.com/wit-switch/assessment-tax/pkg/errorx"
	"go.uber.org/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tax", func() {
	var (
		ctrl    *gomock.Controller
		ctx     context.Context
		service *tax.Services

		mockTaxRepository *mocks.MockTaxRepository
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockTaxRepository = mocks.NewMockTaxRepository(ctrl)
		service = tax.NewService(tax.Dependencies{
			TaxRepository: mockTaxRepository,
		})

		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
		ctrl = nil
	})

	Describe("calculate tax", func() {
		var (
			body domain.TaxCalculate

			mockPersonalDeduct domain.TaxDeduct
			mockTaxDeducts     []domain.TaxDeduct

			mockListTaxDeduct *gomock.Call
		)

		BeforeEach(func() {
			body = domain.TaxCalculate{
				TotalIncome: decimal.NewFromFloat(500000),
				Wht:         decimal.Decimal{},
				Allowances: []domain.Allowance{
					{
						AllowanceType: domain.TaxDeductTypeDonation,
						Amount:        decimal.Decimal{},
					},
				},
			}

			mockPersonalDeduct = domain.TaxDeduct{
				Type:      domain.TaxDeductTypePersonal,
				MinAmount: decimal.NewFromFloat(10000),
				MaxAmount: decimal.NewFromFloat(100000),
				Amount:    decimal.NewFromFloat(60000),
			}
			mockTaxDeducts = []domain.TaxDeduct{
				mockPersonalDeduct,
			}

			mockListTaxDeduct = mockTaxRepository.EXPECT().
				ListTaxDeduct(ctx, domain.GetTaxDeduct{}).
				MinTimes(0)
		})

		When("failed to calculate tax", func() {
			Context("with get list tax deduct error", func() {
				It("should return an error", func() {
					mockListTaxDeduct.Return(nil, errorx.ErrTaxDeductNotFound)

					actual, err := service.Calculate(ctx, body)
					Expect(actual).To(BeNil())
					Expect(err).To(MatchError(errorx.ErrTaxDeductNotFound))
				})
			})

			Context("with get tax deduct from helper error", func() {
				BeforeEach(func() {
					mockTaxDeducts = []domain.TaxDeduct{}
				})
				It("should return an error", func() {
					mockListTaxDeduct.Return(mockTaxDeducts, nil)

					actual, err := service.Calculate(ctx, body)
					Expect(actual).To(BeNil())
					Expect(err).To(MatchError(errorx.ErrTaxDeductNotFound))
				})
			})
		})

		When("success to calculate tax", func() {
			DescribeTable("should return tax without error", func(body domain.TaxCalculate, expected *domain.Tax) {
				mockListTaxDeduct.Return(mockTaxDeducts, nil)

				actual, err := service.Calculate(ctx, body)
				Expect(actual.Tax.InexactFloat64()).To(Equal(expected.Tax.InexactFloat64()))
				Expect(actual.TaxRefund.InexactFloat64()).To(Equal(expected.TaxRefund.InexactFloat64()))
				Expect(err).NotTo(HaveOccurred())
			},
				Entry("with total income is 500000.0",
					domain.TaxCalculate{
						TotalIncome: decimal.NewFromFloat(500000),
						Wht:         decimal.Decimal{},
						Allowances: []domain.Allowance{
							{
								AllowanceType: domain.TaxDeductTypeDonation,
								Amount:        decimal.Decimal{},
							},
						},
					},
					&domain.Tax{
						Tax: decimal.NewFromFloat(29000),
					},
				),
				Entry("with total income is 1000000.0",
					domain.TaxCalculate{
						TotalIncome: decimal.NewFromInt(1000000),
						Wht:         decimal.Decimal{},
						Allowances: []domain.Allowance{
							{
								AllowanceType: domain.TaxDeductTypeDonation,
								Amount:        decimal.Decimal{},
							},
						},
					},
					&domain.Tax{
						Tax: decimal.NewFromFloat(101000),
					},
				),
				Entry("with total income is 500000.0, wth is 25000.0",
					domain.TaxCalculate{
						TotalIncome: decimal.NewFromInt(500000),
						Wht:         decimal.NewFromFloat(25000),
						Allowances: []domain.Allowance{
							{
								AllowanceType: domain.TaxDeductTypeDonation,
								Amount:        decimal.Decimal{},
							},
						},
					},
					&domain.Tax{
						Tax: decimal.NewFromFloat(4000),
					},
				),
				Entry("with total income is 450000.0, wth is 25000.0 and have tax refund",
					domain.TaxCalculate{
						TotalIncome: decimal.NewFromInt(450000),
						Wht:         decimal.NewFromFloat(25000),
						Allowances: []domain.Allowance{
							{
								AllowanceType: domain.TaxDeductTypeDonation,
								Amount:        decimal.Decimal{},
							},
						},
					},
					&domain.Tax{
						Tax:       decimal.NewFromFloat(0),
						TaxRefund: decimal.NewFromFloat(1000),
					},
				),
			)
		})
	})
})
