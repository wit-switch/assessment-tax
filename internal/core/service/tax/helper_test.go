package tax_test

import (
	"github.com/shopspring/decimal"
	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/internal/core/service/tax"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper", func() {

	Describe("get tax deduct by type", func() {
		var (
			taxDeductMap  map[domain.TaxDeductType]domain.TaxDeduct
			taxDeductType domain.TaxDeductType

			mockTaxDeduct *domain.TaxDeduct
		)

		BeforeEach(func() {
			taxDeductMap = map[domain.TaxDeductType]domain.TaxDeduct{
				domain.TaxDeductTypePersonal: {
					Type:      domain.TaxDeductTypePersonal,
					MinAmount: decimal.NewFromFloat(10000),
					MaxAmount: decimal.NewFromFloat(100000),
					Amount:    decimal.NewFromFloat(60000),
				},
				domain.TaxDeductTypeDonation: {},
				domain.TaxDeductTypeKReceipt: {},
			}
			taxDeductType = domain.TaxDeductTypePersonal

			mockTaxDeduct = &domain.TaxDeduct{
				Type:      domain.TaxDeductTypePersonal,
				MinAmount: decimal.NewFromFloat(10000),
				MaxAmount: decimal.NewFromFloat(100000),
				Amount:    decimal.NewFromFloat(60000),
			}
		})

		When("failed to get tax deduct", func() {
			BeforeEach(func() {
				taxDeductType = domain.TaxDeductTypeDonation
			})
			It("should return error", func() {
				actual, err := tax.GetTaxDeductByType(
					taxDeductMap,
					"",
				)

				Expect(actual).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		When("success to get tax deduct", func() {
			It("should return error", func() {
				actual, err := tax.GetTaxDeductByType(
					taxDeductMap,
					taxDeductType,
				)

				Expect(actual).To(Equal(mockTaxDeduct))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
