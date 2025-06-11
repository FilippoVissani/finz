package internal

import (
	"math"
	"testing"
)

func TestCalculateRetirement(t *testing.T) {
	tests := []struct {
		name     string
		input    RetirementInput
		expected RetirementResult
	}{
		{
			name: "Basic retirement calculation",
			input: RetirementInput{
				CurrentAge:          30,
				RetirementAge:       65,
				CurrentSavings:      50000,
				MonthlyContribution: 500,
				WithdrawalRate:      4,
				AnnualYield:         7,
				Inflation:           2,
			},
			expected: RetirementResult{
				CurrentAge:        30,
				RetirementAge:     65,
				YearsToRetirement: 35,
				// These values are calculated by the function
				RetirementSavings:     1475834.89, // This is an approximate value
				AnnualWithdrawal:      59033.40,   // 4% of retirement savings
				MonthlyWithdrawal:     4919.45,    // Annual withdrawal / 12
				RealMonthlyWithdrawal: 2459.86,    // Adjusted for inflation
			},
		},
		{
			name: "Zero yield calculation",
			input: RetirementInput{
				CurrentAge:          40,
				RetirementAge:       60,
				CurrentSavings:      100000,
				MonthlyContribution: 1000,
				WithdrawalRate:      5,
				AnnualYield:         0,
				Inflation:           3,
			},
			expected: RetirementResult{
				CurrentAge:        40,
				RetirementAge:     60,
				YearsToRetirement: 20,
				// With zero yield, we just add up the contributions
				RetirementSavings:     340000,  // 100000 + (1000 * 12 * 20)
				AnnualWithdrawal:      17000,   // 5% of retirement savings
				MonthlyWithdrawal:     1416.67, // Annual withdrawal / 12
				RealMonthlyWithdrawal: 784.73,  // Adjusted for inflation
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateRetirement(tc.input)

			// Check basic fields
			if result.CurrentAge != tc.expected.CurrentAge {
				t.Errorf("CurrentAge = %v, want %v", result.CurrentAge, tc.expected.CurrentAge)
			}
			if result.RetirementAge != tc.expected.RetirementAge {
				t.Errorf("RetirementAge = %v, want %v", result.RetirementAge, tc.expected.RetirementAge)
			}
			if result.YearsToRetirement != tc.expected.YearsToRetirement {
				t.Errorf("YearsToRetirement = %v, want %v", result.YearsToRetirement, tc.expected.YearsToRetirement)
			}

			// For floating point values, use approximate comparison
			// The expected values in the test cases are approximate, so we use a relative tolerance
			const tolerance = 0.01 // 1% tolerance

			if !approximatelyEqual(result.RetirementSavings, tc.expected.RetirementSavings, tolerance) {
				t.Errorf("RetirementSavings = %v, want approximately %v", result.RetirementSavings, tc.expected.RetirementSavings)
			}
			if !approximatelyEqual(result.AnnualWithdrawal, tc.expected.AnnualWithdrawal, tolerance) {
				t.Errorf("AnnualWithdrawal = %v, want approximately %v", result.AnnualWithdrawal, tc.expected.AnnualWithdrawal)
			}
			if !approximatelyEqual(result.MonthlyWithdrawal, tc.expected.MonthlyWithdrawal, tolerance) {
				t.Errorf("MonthlyWithdrawal = %v, want approximately %v", result.MonthlyWithdrawal, tc.expected.MonthlyWithdrawal)
			}
			if !approximatelyEqual(result.RealMonthlyWithdrawal, tc.expected.RealMonthlyWithdrawal, tolerance) {
				t.Errorf("RealMonthlyWithdrawal = %v, want approximately %v", result.RealMonthlyWithdrawal, tc.expected.RealMonthlyWithdrawal)
			}
		})
	}
}

// Helper function to compare floating point values with a relative tolerance
func approximatelyEqual(a, b, tolerance float64) bool {
	if a == b {
		return true
	}

	// Use relative error for non-zero values
	if b != 0 {
		return math.Abs((a-b)/b) < tolerance
	}
	return math.Abs(a-b) < tolerance
}

// TestEdgeCases tests edge cases for the retirement calculator
func TestRetirementEdgeCases(t *testing.T) {
	// Test case where current age equals retirement age
	t.Run("Current age equals retirement age", func(t *testing.T) {
		input := RetirementInput{
			CurrentAge:          65,
			RetirementAge:       65,
			CurrentSavings:      500000,
			MonthlyContribution: 0,
			WithdrawalRate:      4,
			AnnualYield:         5,
			Inflation:           2,
		}

		result := CalculateRetirement(input)

		if result.YearsToRetirement != 0 {
			t.Errorf("YearsToRetirement = %v, want 0", result.YearsToRetirement)
		}

		// The retirement savings should equal the current savings since there's no time for growth
		if result.RetirementSavings != input.CurrentSavings {
			t.Errorf("RetirementSavings = %v, want %v", result.RetirementSavings, input.CurrentSavings)
		}
	})

	// Test case with negative yield (market downturn)
	t.Run("Negative yield", func(t *testing.T) {
		input := RetirementInput{
			CurrentAge:          40,
			RetirementAge:       50,
			CurrentSavings:      200000,
			MonthlyContribution: 1000,
			WithdrawalRate:      3,
			AnnualYield:         -2,
			Inflation:           1,
		}

		result := CalculateRetirement(input)

		// Ensure the calculation completes without errors
		// The retirement savings should be less than the sum of current savings and contributions
		totalContributions := input.CurrentSavings + (float64(result.YearsToRetirement) * 12 * input.MonthlyContribution)
		if result.RetirementSavings > totalContributions {
			t.Errorf("With negative yield, RetirementSavings (%v) should be less than total contributions (%v)",
				result.RetirementSavings, totalContributions)
		}
	})
}
