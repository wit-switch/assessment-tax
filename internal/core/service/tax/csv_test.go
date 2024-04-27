package tax_test

import (
	"bytes"
	"context"
	"encoding/csv"

	"github.com/gocarina/gocsv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/wit-switch/assessment-tax/internal/core/domain"
	"github.com/wit-switch/assessment-tax/internal/core/service/tax"
	"github.com/wit-switch/assessment-tax/mocks"
	"github.com/wit-switch/assessment-tax/pkg/errorx"
	"go.uber.org/mock/gomock"
)

var _ = Describe("CSV", func() {
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

	Describe("calculate tax from csv", func() {
		var (
			file csv.Reader

			mockPersonalDeduct domain.TaxDeduct
			mockDonationDeduct domain.TaxDeduct
			mockTaxDeducts     []domain.TaxDeduct

			mockListTaxDeduct *gomock.Call

			zero decimal.Decimal
		)

		BeforeEach(func() {
			zero = decimal.NewFromFloat(0)

			body := new(bytes.Buffer)
			writer := csv.NewWriter(body)
			_ = writer.Write([]string{"totalIncome", "wht", "donation"})
			_ = writer.Write([]string{"500000", "0", "0"})
			_ = writer.Write([]string{"600000", "40000", "20000"})
			_ = writer.Write([]string{"750000", "50000", "15000"})
			writer.Flush()

			r := csv.NewReader(body)

			file = *r

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
			mockTaxDeducts = []domain.TaxDeduct{
				mockPersonalDeduct,
				mockDonationDeduct,
			}

			mockListTaxDeduct = mockTaxRepository.EXPECT().
				ListTaxDeduct(ctx, domain.GetTaxDeduct{}).
				MinTimes(0)
		})

		When("failed to calculate tax from csv", func() {
			Context("with validate file rror", func() {
				BeforeEach(func() {
					body := new(bytes.Buffer)
					writer := csv.NewWriter(body)
					writer.Flush()

					r := csv.NewReader(body)

					file = *r
				})
				It("should return error validate file", func() {
					actual, err := service.CalculateFromCSV(ctx, file)
					Expect(actual).To(BeNil())
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(gocsv.ErrEmptyCSVFile))
				})
			})

			Context("with totalIncome should greater than or equal 0", func() {
				BeforeEach(func() {
					body := new(bytes.Buffer)
					writer := csv.NewWriter(body)
					_ = writer.Write([]string{"totalIncome", "wht", "donation"})
					_ = writer.Write([]string{"-500000", "10", "0"})
					writer.Flush()

					r := csv.NewReader(body)

					file = *r
				})
				It("should return error validate fail", func() {
					actual, err := service.CalculateFromCSV(ctx, file)
					Expect(actual).To(BeNil())
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(errorx.ErrValidationFail))
				})
			})

			Context("with wht should greater than or equal 0", func() {
				BeforeEach(func() {
					body := new(bytes.Buffer)
					writer := csv.NewWriter(body)
					_ = writer.Write([]string{"totalIncome", "wht", "donation"})
					_ = writer.Write([]string{"500000", "-10", "0"})
					writer.Flush()

					r := csv.NewReader(body)

					file = *r
				})
				It("should return error validate fail", func() {
					actual, err := service.CalculateFromCSV(ctx, file)
					Expect(actual).To(BeNil())
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(errorx.ErrValidationFail))
				})
			})

			Context("with wht should less than or equal totalIncome", func() {
				BeforeEach(func() {
					body := new(bytes.Buffer)
					writer := csv.NewWriter(body)
					_ = writer.Write([]string{"totalIncome", "wht", "donation"})
					_ = writer.Write([]string{"10", "13", "0"})
					writer.Flush()

					r := csv.NewReader(body)

					file = *r
				})
				It("should return error validate fail", func() {
					actual, err := service.CalculateFromCSV(ctx, file)
					Expect(actual).To(BeNil())
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(errorx.ErrValidationFail))
				})
			})
		})

		When("success to calculate tax from csv", func() {
			It("should return error validate fail", func() {
				mockListTaxDeduct.Return(mockTaxDeducts, nil)

				actual, err := service.CalculateFromCSV(ctx, file)

				expected := []domain.TaxCSV{
					{
						TotalIncome: decimal.NewFromFloat(500000),
						Tax:         decimal.NewFromFloat(29000),
						TaxRefund:   zero,
					},
					{
						TotalIncome: decimal.NewFromFloat(600000),
						Tax:         zero,
						TaxRefund:   decimal.NewFromFloat(2000),
					},
					{
						TotalIncome: decimal.NewFromFloat(750000),
						Tax:         decimal.NewFromFloat(11250),
						TaxRefund:   zero,
					},
				}

				Expect(actual).To(BeComparableTo(expected))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
