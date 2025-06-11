package internal

import (
	"math"
)

// RetirementInput represents the input parameters for retirement calculation
type RetirementInput struct {
	CurrentAge          int
	RetirementAge       int
	CurrentSavings      float64
	MonthlyContribution float64
	WithdrawalRate      float64
	AnnualYield         float64
	Inflation           float64
}

// RetirementResult represents the output of retirement calculation
type RetirementResult struct {
	CurrentAge            int
	RetirementAge         int
	YearsToRetirement     int
	RetirementSavings     float64
	AnnualWithdrawal      float64
	MonthlyWithdrawal     float64
	RealMonthlyWithdrawal float64
}

func CalculateRetirement(input RetirementInput) RetirementResult {
	yearsToRetirement := input.RetirementAge - input.CurrentAge
	monthsToRetirement := yearsToRetirement * 12

	// Calculate retirement savings at retirement age
	monthlyRate := input.AnnualYield / 100 / 12

	// Future value calculation
	retirementSavings := input.CurrentSavings * math.Pow(1+monthlyRate, float64(monthsToRetirement))
	if monthlyRate > 0 {
		retirementSavings += input.MonthlyContribution * (math.Pow(1+monthlyRate, float64(monthsToRetirement)) - 1) / monthlyRate
	} else {
		retirementSavings += input.MonthlyContribution * float64(monthsToRetirement)
	}

	// Calculate annual withdrawal amount
	annualWithdrawal := retirementSavings * (input.WithdrawalRate / 100)
	monthlyWithdrawal := annualWithdrawal / 12

	// Adjust for inflation
	realMonthlyWithdrawal := monthlyWithdrawal / math.Pow(1+input.Inflation/100, float64(yearsToRetirement))

	return RetirementResult{
		CurrentAge:            input.CurrentAge,
		RetirementAge:         input.RetirementAge,
		YearsToRetirement:     yearsToRetirement,
		RetirementSavings:     retirementSavings,
		AnnualWithdrawal:      annualWithdrawal,
		MonthlyWithdrawal:     monthlyWithdrawal,
		RealMonthlyWithdrawal: realMonthlyWithdrawal,
	}
}
