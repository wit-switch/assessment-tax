package tax_test

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/internal/core/service/tax"
	"github.com/wit-switch/assessment-tax/pkg/errorx"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper", func() {

	Describe("get tax deduct by type", func() {
		var (
			taxDeductMap  map[domain.TaxDeductType]domain.TaxDeduct
			taxDeductType domain.TaxDeductType
		)

		BeforeEach(func() {
			taxDeductMap = map[domain.TaxDeductType]domain.TaxDeduct{
				domain.TaxDeductTypePersonal: {
					Type:      domain.TaxDeductTypePersonal,
					MinAmount: decimal.NewFromFloat(10000),
					MaxAmount: decimal.NewFromFloat(100000),
					Amount:    decimal.NewFromFloat(60000),
				},
				domain.TaxDeductTypeDonation: {
					Type:      domain.TaxDeductTypeDonation,
					MinAmount: decimal.NewFromFloat(0),
					MaxAmount: decimal.NewFromFloat(100000),
					Amount:    decimal.NewFromFloat(100000),
				},
				domain.TaxDeductTypeKReceipt: {},
			}
			taxDeductType = domain.TaxDeductTypePersonal
		})

		When("failed to get tax deduct", func() {
			BeforeEach(func() {
				taxDeductType = "xxx"
			})
			It("should return error", func() {
				actual, err := tax.GetTaxDeductByType(
					taxDeductMap,
					"",
				)

				Expect(actual).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errorx.ErrTaxDeductNotFound))
			})
		})

		When("success to get tax deduct", func() {
			Context("with donation type", func() {
				BeforeEach(func() {
					taxDeductType = domain.TaxDeductTypeDonation
				})
				It("should return tax deduct", func() {
					actual, err := tax.GetTaxDeductByType(
						taxDeductMap,
						taxDeductType,
					)

					expected := &domain.TaxDeduct{
						Type:      domain.TaxDeductTypeDonation,
						MinAmount: decimal.NewFromFloat(0),
						MaxAmount: decimal.NewFromFloat(100000),
						Amount:    decimal.NewFromFloat(100000),
					}

					Expect(actual.Type).To(Equal(expected.Type))
					Expect(actual.MinAmount.InexactFloat64()).To(Equal(expected.MinAmount.InexactFloat64()))
					Expect(actual.MaxAmount.InexactFloat64()).To(Equal(expected.MaxAmount.InexactFloat64()))
					Expect(actual.Amount.InexactFloat64()).To(Equal(expected.Amount.InexactFloat64()))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			It("should return tax deduct", func() {
				actual, err := tax.GetTaxDeductByType(
					taxDeductMap,
					taxDeductType,
				)

				expected := &domain.TaxDeduct{
					Type:      domain.TaxDeductTypePersonal,
					MinAmount: decimal.NewFromFloat(10000),
					MaxAmount: decimal.NewFromFloat(100000),
					Amount:    decimal.NewFromFloat(60000),
				}

				Expect(actual.Type).To(Equal(expected.Type))
				Expect(actual.MinAmount.InexactFloat64()).To(Equal(expected.MinAmount.InexactFloat64()))
				Expect(actual.MaxAmount.InexactFloat64()).To(Equal(expected.MaxAmount.InexactFloat64()))
				Expect(actual.Amount.InexactFloat64()).To(Equal(expected.Amount.InexactFloat64()))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("get allowance amount", func() {
		var (
			ctx             context.Context
			taxDeduct       *domain.TaxDeduct
			allowanceAmount decimal.Decimal
		)

		BeforeEach(func() {
			ctx = context.Background()
			taxDeduct = &domain.TaxDeduct{
				Type:      domain.TaxDeductTypeDonation,
				MinAmount: decimal.NewFromFloat(0),
				MaxAmount: decimal.NewFromFloat(100000),
				Amount:    decimal.NewFromFloat(100000),
			}
			allowanceAmount = decimal.NewFromFloat(10000)
		})

		When("more than limit", func() {
			BeforeEach(func() {
				allowanceAmount = decimal.NewFromFloat(200000)
			})
			It("should return limit allowance amount", func() {
				actual, err := tax.GetAllowanceAmount(ctx,
					taxDeduct,
					allowanceAmount,
				)

				expected := decimal.NewFromFloat(100000)

				Expect(actual.InexactFloat64()).To(Equal(expected.InexactFloat64()))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("less than limit", func() {
			BeforeEach(func() {
				allowanceAmount = decimal.NewFromFloat(-200000)
			})
			It("should return limit allowance amount", func() {
				actual, err := tax.GetAllowanceAmount(ctx,
					taxDeduct,
					allowanceAmount,
				)

				Expect(actual).NotTo(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errorx.ErrAmountLessThanLimit))
			})
		})

		When("normal case", func() {
			It("should return allowance amount", func() {
				actual, err := tax.GetAllowanceAmount(ctx,
					taxDeduct,
					allowanceAmount,
				)

				expected := decimal.NewFromFloat(10000)

				Expect(actual.InexactFloat64()).To(Equal(expected.InexactFloat64()))
				Expect(err).NotTo(HaveOccurred())
			})
		})

	})
})
