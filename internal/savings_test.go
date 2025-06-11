package internal

import (
	"testing"
)

func TestCalculateSavings(t *testing.T) {
	tests := []struct {
		name     string
		input    SavingsInput
		expected SavingsResult
	}{
		{
			name: "Basic savings calculation",
			input: SavingsInput{
				Initial:        1000,
				MonthlyDeposit: 100,
				AnnualYield:    5,
				Inflation:      2,
				Years:          10,
			},
			expected: SavingsResult{
				Initial:         1000,
				MonthlyDeposit:  100,
				FutureValue:     17175.24, // Calculated with compound interest formula
				RealFutureValue: 14093.64, // FutureValue adjusted for inflation
				TotalDeposits:   13000,    // 1000 + (100 * 12 * 10)
				InterestEarned:  4175.24,  // FutureValue - TotalDeposits
				Years:           10,
				NumberOfMonths:  120,
			},
		},
		{
			name: "Zero yield savings",
			input: SavingsInput{
				Initial:        5000,
				MonthlyDeposit: 200,
				AnnualYield:    0,
				Inflation:      2,
				Years:          5,
			},
			expected: SavingsResult{
				Initial:         5000,
				MonthlyDeposit:  200,
				FutureValue:     17000,    // 5000 + (200 * 12 * 5)
				RealFutureValue: 15401.33, // FutureValue adjusted for inflation
				TotalDeposits:   17000,    // Same as future value with zero yield
				InterestEarned:  0,        // No interest with zero yield
				Years:           5,
				NumberOfMonths:  60,
			},
		},
		{
			name: "High yield, long term savings",
			input: SavingsInput{
				Initial:        2000,
				MonthlyDeposit: 500,
				AnnualYield:    8,
				Inflation:      2.5,
				Years:          20,
			},
			expected: SavingsResult{
				Initial:         2000,
				MonthlyDeposit:  500,
				FutureValue:     304363.81, // Calculated with compound interest formula
				RealFutureValue: 186147.17, // FutureValue adjusted for inflation
				TotalDeposits:   122000,    // 2000 + (500 * 12 * 20)
				InterestEarned:  182363.81, // FutureValue - TotalDeposits
				Years:           20,
				NumberOfMonths:  240,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateSavings(tc.input)

			// Check basic fields
			if result.Initial != tc.expected.Initial {
				t.Errorf("Initial = %v, want %v", result.Initial, tc.expected.Initial)
			}
			if result.MonthlyDeposit != tc.expected.MonthlyDeposit {
				t.Errorf("MonthlyDeposit = %v, want %v", result.MonthlyDeposit, tc.expected.MonthlyDeposit)
			}
			if result.Years != tc.expected.Years {
				t.Errorf("Years = %v, want %v", result.Years, tc.expected.Years)
			}
			if result.NumberOfMonths != tc.expected.NumberOfMonths {
				t.Errorf("NumberOfMonths = %v, want %v", result.NumberOfMonths, tc.expected.NumberOfMonths)
			}

			// For floating point values, use approximate comparison
			const tolerance = 0.01 // 1% tolerance

			if !approximatelyEqual(result.FutureValue, tc.expected.FutureValue, tolerance) {
				t.Errorf("FutureValue = %v, want approximately %v", result.FutureValue, tc.expected.FutureValue)
			}
			if !approximatelyEqual(result.RealFutureValue, tc.expected.RealFutureValue, tolerance) {
				t.Errorf("RealFutureValue = %v, want approximately %v", result.RealFutureValue, tc.expected.RealFutureValue)
			}
			if !approximatelyEqual(result.TotalDeposits, tc.expected.TotalDeposits, tolerance) {
				t.Errorf("TotalDeposits = %v, want approximately %v", result.TotalDeposits, tc.expected.TotalDeposits)
			}
			if !approximatelyEqual(result.InterestEarned, tc.expected.InterestEarned, tolerance) {
				t.Errorf("InterestEarned = %v, want approximately %v", result.InterestEarned, tc.expected.InterestEarned)
			}
		})
	}
}

// TestSavingsEdgeCases tests edge cases for the savings calculator
func TestSavingsEdgeCases(t *testing.T) {
	// Test case with zero initial amount
	t.Run("Zero initial amount", func(t *testing.T) {
		input := SavingsInput{
			Initial:        0,
			MonthlyDeposit: 100,
			AnnualYield:    6,
			Inflation:      2,
			Years:          5,
		}

		result := CalculateSavings(input)

		// Future value should only reflect the monthly deposits plus interest
		if result.FutureValue <= result.TotalDeposits && input.AnnualYield > 0 {
			t.Errorf("With positive yield, FutureValue (%v) should be greater than TotalDeposits (%v)",
				result.FutureValue, result.TotalDeposits)
		}

		// Total deposits should be monthly deposit * number of months
		expectedTotalDeposits := input.MonthlyDeposit * float64(result.NumberOfMonths)
		if !approximatelyEqual(result.TotalDeposits, expectedTotalDeposits, 0.001) {
			t.Errorf("TotalDeposits = %v, want %v", result.TotalDeposits, expectedTotalDeposits)
		}
	})

	// Test case with zero monthly deposit
	t.Run("Zero monthly deposit", func(t *testing.T) {
		input := SavingsInput{
			Initial:        10000,
			MonthlyDeposit: 0,
			AnnualYield:    4,
			Inflation:      2,
			Years:          10,
		}

		result := CalculateSavings(input)

		// Future value should only reflect the initial amount plus interest
		// This is a rough approximation, so we use a larger tolerance
		if result.FutureValue < input.Initial && input.AnnualYield > 0 {
			t.Errorf("With positive yield, FutureValue (%v) should be greater than Initial (%v)",
				result.FutureValue, input.Initial)
		}

		// Total deposits should equal initial amount
		if result.TotalDeposits != input.Initial {
			t.Errorf("TotalDeposits = %v, want %v", result.TotalDeposits, input.Initial)
		}
	})

	// Test case with zero years
	t.Run("Zero years", func(t *testing.T) {
		input := SavingsInput{
			Initial:        5000,
			MonthlyDeposit: 200,
			AnnualYield:    7,
			Inflation:      2,
			Years:          0,
		}

		result := CalculateSavings(input)

		// With zero years, future value should equal initial amount
		if result.FutureValue != input.Initial {
			t.Errorf("FutureValue = %v, want %v", result.FutureValue, input.Initial)
		}

		// Total deposits should equal initial amount
		if result.TotalDeposits != input.Initial {
			t.Errorf("TotalDeposits = %v, want %v", result.TotalDeposits, input.Initial)
		}

		// Interest earned should be zero
		if result.InterestEarned != 0 {
			t.Errorf("InterestEarned = %v, want 0", result.InterestEarned)
		}

		// Number of months should be zero
		if result.NumberOfMonths != 0 {
			t.Errorf("NumberOfMonths = %v, want 0", result.NumberOfMonths)
		}
	})

	// Test case with negative yield (market downturn)
	t.Run("Negative yield", func(t *testing.T) {
		input := SavingsInput{
			Initial:        20000,
			MonthlyDeposit: 300,
			AnnualYield:    -3,
			Inflation:      2,
			Years:          2,
		}

		result := CalculateSavings(input)

		// With negative yield, future value should be less than total deposits
		if result.FutureValue >= result.TotalDeposits {
			t.Errorf("With negative yield, FutureValue (%v) should be less than TotalDeposits (%v)",
				result.FutureValue, result.TotalDeposits)
		}

		// Interest earned should be negative
		if result.InterestEarned >= 0 {
			t.Errorf("With negative yield, InterestEarned (%v) should be negative", result.InterestEarned)
		}
	})
}
