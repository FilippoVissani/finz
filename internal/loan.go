package internal

import (
	"math"
)

// LoanInput represents the input parameters for loan calculation
type LoanInput struct {
	Principal float64
	Rate      float64
	Years     int
	Monthly   bool
}

// MonthlyBreakdown represents a single month's payment breakdown
type MonthlyBreakdown struct {
	Month            int
	Payment          float64
	PrincipalPayment float64
	InterestPayment  float64
	RemainingBalance float64
}

// LoanResult represents the output of loan calculation
type LoanResult struct {
	Principal        float64
	MonthlyPayment   float64
	TotalPaid        float64
	TotalInterest    float64
	Years            int
	NumberOfPayments int
	MonthlyDetails   []MonthlyBreakdown
}

func CalculateLoan(input LoanInput) LoanResult {
	monthlyRate := input.Rate / 100 / 12
	numberOfPayments := input.Years * 12

	// Monthly payment formula: P * r * (1+r)^n / ((1+r)^n - 1)
	monthlyPayment := input.Principal * monthlyRate * math.Pow(1+monthlyRate, float64(numberOfPayments)) /
		(math.Pow(1+monthlyRate, float64(numberOfPayments)) - 1)

	totalPaid := monthlyPayment * float64(numberOfPayments)
	totalInterest := totalPaid - input.Principal

	result := LoanResult{
		Principal:        input.Principal,
		MonthlyPayment:   monthlyPayment,
		TotalPaid:        totalPaid,
		TotalInterest:    totalInterest,
		Years:            input.Years,
		NumberOfPayments: numberOfPayments,
		MonthlyDetails:   []MonthlyBreakdown{},
	}

	if input.Monthly {
		balance := input.Principal
		for i := 1; i <= 12; i++ { // Show first year only
			interestPayment := balance * monthlyRate
			principalPayment := monthlyPayment - interestPayment
			balance -= principalPayment

			result.MonthlyDetails = append(result.MonthlyDetails, MonthlyBreakdown{
				Month:            i,
				Payment:          monthlyPayment,
				PrincipalPayment: principalPayment,
				InterestPayment:  interestPayment,
				RemainingBalance: balance,
			})
		}
	}

	return result
}
