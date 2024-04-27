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
			mockDonationDeduct domain.TaxDeduct
			mockKReceiptDeduct domain.TaxDeduct
			mockTaxDeducts     []domain.TaxDeduct

			mockListTaxDeduct *gomock.Call

			zero decimal.Decimal
		)

		BeforeEach(func() {
			zero = decimal.NewFromFloat(0)

			body = domain.TaxCalculate{
				TotalIncome: decimal.NewFromFloat(500000),
				Wht:         zero,
				Allowances: []domain.Allowance{
					{
						AllowanceType: domain.TaxDeductTypeDonation,
						Amount:        decimal.NewFromFloat(200000),
					},
					{
						AllowanceType: domain.TaxDeductTypeKReceipt,
						Amount:        decimal.NewFromFloat(1),
					},
				},
			}

			mockPersonalDeduct = domain.TaxDeduct{
				Type:      domain.TaxDeductTypePersonal,
				MinAmount: decimal.NewFromFloat(10000),
				MaxAmount: decimal.NewFromFloat(100000),
				Amount:    decimal.NewFromFloat(60000),
			}
			mockDonationDeduct = domain.TaxDeduct{
				Type:      domain.TaxDeductTypeDonation,
				MinAmount: zero,
				MaxAmount: decimal.NewFromFloat(100000),
				Amount:    decimal.NewFromFloat(100000),
			}
			mockKReceiptDeduct = domain.TaxDeduct{
				Type:      domain.TaxDeductTypeKReceipt,
				MinAmount: decimal.NewFromFloat(1),
				MaxAmount: decimal.NewFromFloat(100000),
				Amount:    decimal.NewFromFloat(50000),
			}
			mockTaxDeducts = []domain.TaxDeduct{
				mockPersonalDeduct,
				mockDonationDeduct,
				mockKReceiptDeduct,
			}

			mockListTaxDeduct = mockTaxRepository.EXPECT().
				ListTaxDeduct(ctx, domain.GetTaxDeduct{}).
				MinTimes(0)
		})

		When("failed to calculate tax", func() {
			Context("with get list tax deduct error", func() {
				It("should return an error", func() {
					mockListTaxDeduct.Return(nil, errorx.ErrTaxDeductNotFound)

					actual, err := service.Calculate(ctx, body, true)
					Expect(actual).To(BeNil())
					Expect(err).To(MatchError(errorx.ErrTaxDeductNotFound))
				})
			})

			Context("with get personal tax deduct from helper error", func() {
				BeforeEach(func() {
					mockTaxDeducts = []domain.TaxDeduct{}
				})
				It("should return an error", func() {
					mockListTaxDeduct.Return(mockTaxDeducts, nil)

					actual, err := service.Calculate(ctx, body, true)
					Expect(actual).To(BeNil())
					Expect(err).To(MatchError(errorx.ErrTaxDeductNotFound))
				})
			})

			Context("with get donation tax deduct from helper error", func() {
				BeforeEach(func() {
					mockTaxDeducts = []domain.TaxDeduct{
						mockPersonalDeduct,
					}
				})
				It("should return an error", func() {
					mockListTaxDeduct.Return(mockTaxDeducts, nil)

					actual, err := service.Calculate(ctx, body, true)
					Expect(actual).To(BeNil())
					Expect(err).To(MatchError(errorx.ErrTaxDeductNotFound))
				})
			})

			Context("with get donation allowance amount from helper error", func() {
				BeforeEach(func() {
					mockDonationDeduct.MinAmount = decimal.NewFromFloat(1000000)
					mockTaxDeducts = []domain.TaxDeduct{
						mockPersonalDeduct,
						mockDonationDeduct,
						mockKReceiptDeduct,
					}
				})
				It("should return an error", func() {
					mockListTaxDeduct.Return(mockTaxDeducts, nil)

					actual, err := service.Calculate(ctx, body, true)
					Expect(actual).To(BeNil())
					Expect(err).To(MatchError(errorx.ErrAmountLessThanLimit))
				})
			})
		})

		When("success to calculate tax", func() {
			DescribeTable("should return tax without error", func(body domain.TaxCalculate, expected *domain.Tax) {
				mockListTaxDeduct.Return(mockTaxDeducts, nil)

				actual, err := service.Calculate(ctx, body, true)
				Expect(actual.Tax.InexactFloat64()).To(Equal(expected.Tax.InexactFloat64()))
				Expect(actual.TaxRefund.InexactFloat64()).To(Equal(expected.TaxRefund.InexactFloat64()))
				Expect(actual.TaxLevel).To(BeComparableTo(expected.TaxLevel))
				Expect(err).NotTo(HaveOccurred())
			},
				Entry("with total income is 500000.0",
					domain.TaxCalculate{
						TotalIncome: decimal.NewFromFloat(500000),
						Wht:         zero,
						Allowances: []domain.Allowance{
							{
								AllowanceType: domain.TaxDeductTypeDonation,
								Amount:        zero,
							},
							{
								AllowanceType: domain.TaxDeductTypeKReceipt,
								Amount:        decimal.NewFromFloat(1),
							},
						},
					},
					&domain.Tax{
						Tax: decimal.NewFromFloat(28999.9),
						TaxLevel: []domain.TaxLevel{
							{
								Level: "0-150,000",
								Tax:   zero,
							},
							{
								Level: "150,001-500,000",
								Tax:   decimal.NewFromFloat(28999.9),
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
					},
				),
				Entry("with total income is 1000000.0",
					domain.TaxCalculate{
						TotalIncome: decimal.NewFromInt(1000000),
						Wht:         zero,
						Allowances: []domain.Allowance{
							{
								AllowanceType: domain.TaxDeductTypeDonation,
								Amount:        zero,
							},
							{
								AllowanceType: domain.TaxDeductTypeKReceipt,
								Amount:        decimal.NewFromFloat(1),
							},
						},
					},
					&domain.Tax{
						Tax: decimal.NewFromFloat(100999.85),
						TaxLevel: []domain.TaxLevel{
							{
								Level: "0-150,000",
								Tax:   zero,
							},
							{
								Level: "150,001-500,000",
								Tax:   decimal.NewFromFloat(35000),
							},
							{
								Level: "500,001-1,000,000",
								Tax:   decimal.NewFromFloat(65999.85),
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
					},
				),
				Entry("with total income is 500000.0, wth is 25000.0",
					domain.TaxCalculate{
						TotalIncome: decimal.NewFromInt(500000),
						Wht:         decimal.NewFromFloat(25000),
						Allowances: []domain.Allowance{
							{
								AllowanceType: domain.TaxDeductTypeDonation,
								Amount:        zero,
							},
							{
								AllowanceType: domain.TaxDeductTypeKReceipt,
								Amount:        decimal.NewFromFloat(1),
							},
						},
					},
					&domain.Tax{
						Tax: decimal.NewFromFloat(3999.9),
						TaxLevel: []domain.TaxLevel{
							{
								Level: "0-150,000",
								Tax:   zero,
							},
							{
								Level: "150,001-500,000",
								Tax:   decimal.NewFromFloat(28999.9),
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
					},
				),
				Entry("with total income is 450000.0, wth is 25000.0 and have tax refund",
					domain.TaxCalculate{
						TotalIncome: decimal.NewFromInt(450000),
						Wht:         decimal.NewFromFloat(25000),
						Allowances: []domain.Allowance{
							{
								AllowanceType: domain.TaxDeductTypeDonation,
								Amount:        zero,
							},
							{
								AllowanceType: domain.TaxDeductTypeKReceipt,
								Amount:        decimal.NewFromFloat(1),
							},
						},
					},
					&domain.Tax{
						Tax:       decimal.NewFromFloat(0),
						TaxRefund: decimal.NewFromFloat(1000.1),
						TaxLevel: []domain.TaxLevel{
							{
								Level: "0-150,000",
								Tax:   zero,
							},
							{
								Level: "150,001-500,000",
								Tax:   decimal.NewFromFloat(23999.9),
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
					},
				),
			)
		})
	})
})
