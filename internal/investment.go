package internal

import (
	"math"
)

// InvestmentInput represents the input parameters for investment calculation
type InvestmentInput struct {
	Principal   float64
	AnnualYield float64
	TaxRate     float64
	Inflation   float64
	Years       int
}

// InvestmentResult represents the output of investment calculation
type InvestmentResult struct {
	Principal      float64
	NetFutureValue float64
	RealValue      float64
	TaxPaid        float64
	Years          int
}

func CalculateInvestment(input InvestmentInput) InvestmentResult {
	rate := input.AnnualYield / 100
	tax := input.TaxRate / 100
	inf := input.Inflation / 100

	futureValue := input.Principal * math.Pow(1+rate, float64(input.Years))
	profit := futureValue - input.Principal
	taxPaid := profit * tax
	netFutureValue := futureValue - taxPaid

	// Adjust for inflation
	realValue := netFutureValue / math.Pow(1+inf, float64(input.Years))

	return InvestmentResult{
		Principal:      input.Principal,
		NetFutureValue: netFutureValue,
		RealValue:      realValue,
		TaxPaid:        taxPaid,
		Years:          input.Years,
	}
}
