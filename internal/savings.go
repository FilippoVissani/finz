package internal

import (
	"math"
)

// SavingsInput represents the input parameters for savings calculation
type SavingsInput struct {
	Initial        float64
	MonthlyDeposit float64
	AnnualYield    float64
	Years          int
	Inflation      float64
}

// SavingsResult represents the output of savings calculation
type SavingsResult struct {
	Initial         float64
	MonthlyDeposit  float64
	FutureValue     float64
	RealFutureValue float64
	TotalDeposits   float64
	InterestEarned  float64
	Years           int
	NumberOfMonths  int
}

func CalculateSavings(input SavingsInput) SavingsResult {
	monthlyRate := input.AnnualYield / 100 / 12
	numberOfMonths := input.Years * 12

	// Calculate future value with regular deposits
	// FV = P(1+r)^n + PMT * ((1+r)^n - 1) / r
	futureValue := input.Initial * math.Pow(1+monthlyRate, float64(numberOfMonths))
	if monthlyRate > 0 {
		futureValue += input.MonthlyDeposit * (math.Pow(1+monthlyRate, float64(numberOfMonths)) - 1) / monthlyRate
	} else {
		futureValue += input.MonthlyDeposit * float64(numberOfMonths)
	}

	totalDeposits := input.Initial + (input.MonthlyDeposit * float64(numberOfMonths))
	interestEarned := futureValue - totalDeposits

	// Adjust future value for inflation
	realFutureValue := futureValue
	if input.Inflation > 0 {
		realFutureValue = futureValue / math.Pow(1+input.Inflation/100, float64(input.Years))
	}

	return SavingsResult{
		Initial:         input.Initial,
		MonthlyDeposit:  input.MonthlyDeposit,
		FutureValue:     futureValue,
		RealFutureValue: realFutureValue,
		TotalDeposits:   totalDeposits,
		InterestEarned:  interestEarned,
		Years:           input.Years,
		NumberOfMonths:  numberOfMonths,
	}
}
