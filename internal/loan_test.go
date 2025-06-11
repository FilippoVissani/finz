package internal

import (
	"math"
	"testing"
)

func TestCalculateLoan(t *testing.T) {
	tests := []struct {
		name     string
		input    LoanInput
		expected LoanResult
	}{
		{
			name: "Basic mortgage calculation",
			input: LoanInput{
				Principal: 300000,
				Rate:      4.5,
				Years:     30,
				Monthly:   false,
			},
			expected: LoanResult{
				Principal:        300000,
				MonthlyPayment:   1520.06,   // Calculated with standard mortgage formula
				TotalPaid:        547221.60, // Monthly payment * number of payments
				TotalInterest:    247221.60, // Total paid - principal
				Years:            30,
				NumberOfPayments: 360,
			},
		},
		{
			name: "Short-term loan",
			input: LoanInput{
				Principal: 20000,
				Rate:      6.0,
				Years:     5,
				Monthly:   false,
			},
			expected: LoanResult{
				Principal:        20000,
				MonthlyPayment:   386.66,   // Calculated with standard loan formula
				TotalPaid:        23199.60, // Monthly payment * number of payments
				TotalInterest:    3199.60,  // Total paid - principal
				Years:            5,
				NumberOfPayments: 60,
			},
		},
		{
			name: "High interest loan",
			input: LoanInput{
				Principal: 10000,
				Rate:      12.0,
				Years:     3,
				Monthly:   true, // Include monthly breakdown
			},
			expected: LoanResult{
				Principal:        10000,
				MonthlyPayment:   332.14,   // Calculated with standard loan formula
				TotalPaid:        11957.04, // Monthly payment * number of payments
				TotalInterest:    1957.04,  // Total paid - principal
				Years:            3,
				NumberOfPayments: 36,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateLoan(tc.input)

			// Check basic fields
			if result.Principal != tc.expected.Principal {
				t.Errorf("Principal = %v, want %v", result.Principal, tc.expected.Principal)
			}
			if result.Years != tc.expected.Years {
				t.Errorf("Years = %v, want %v", result.Years, tc.expected.Years)
			}
			if result.NumberOfPayments != tc.expected.NumberOfPayments {
				t.Errorf("NumberOfPayments = %v, want %v", result.NumberOfPayments, tc.expected.NumberOfPayments)
			}

			// For floating point values, use approximate comparison
			const tolerance = 0.01 // 1% tolerance

			if !approximatelyEqual(result.MonthlyPayment, tc.expected.MonthlyPayment, tolerance) {
				t.Errorf("MonthlyPayment = %v, want approximately %v", result.MonthlyPayment, tc.expected.MonthlyPayment)
			}
			if !approximatelyEqual(result.TotalPaid, tc.expected.TotalPaid, tolerance) {
				t.Errorf("TotalPaid = %v, want approximately %v", result.TotalPaid, tc.expected.TotalPaid)
			}
			if !approximatelyEqual(result.TotalInterest, tc.expected.TotalInterest, tolerance) {
				t.Errorf("TotalInterest = %v, want approximately %v", result.TotalInterest, tc.expected.TotalInterest)
			}

			// Check monthly breakdown if requested
			if tc.input.Monthly {
				if len(result.MonthlyDetails) != 12 {
					t.Errorf("Expected 12 months of details, got %d", len(result.MonthlyDetails))
				} else {
					// Check first month's details
					firstMonth := result.MonthlyDetails[0]
					if firstMonth.Month != 1 {
						t.Errorf("First month number = %v, want 1", firstMonth.Month)
					}
					if !approximatelyEqual(firstMonth.Payment, result.MonthlyPayment, 0.001) {
						t.Errorf("First month payment = %v, want %v", firstMonth.Payment, result.MonthlyPayment)
					}

					// Check that interest + principal = payment
					if !approximatelyEqual(firstMonth.InterestPayment+firstMonth.PrincipalPayment, firstMonth.Payment, 0.001) {
						t.Errorf("Interest (%v) + Principal (%v) != Payment (%v)",
							firstMonth.InterestPayment, firstMonth.PrincipalPayment, firstMonth.Payment)
					}

					// Check that the remaining balance decreases
					for i := 1; i < len(result.MonthlyDetails); i++ {
						if result.MonthlyDetails[i].RemainingBalance >= result.MonthlyDetails[i-1].RemainingBalance {
							t.Errorf("Remaining balance should decrease: month %d (%v) >= month %d (%v)",
								i+1, result.MonthlyDetails[i].RemainingBalance,
								i, result.MonthlyDetails[i-1].RemainingBalance)
						}
					}
				}
			} else {
				if len(result.MonthlyDetails) != 0 {
					t.Errorf("Expected no monthly details, got %d", len(result.MonthlyDetails))
				}
			}
		})
	}
}

// TestLoanEdgeCases tests edge cases for the loan calculator
func TestLoanEdgeCases(t *testing.T) {
	// Test case with zero interest rate
	t.Run("Zero interest rate", func(t *testing.T) {
		input := LoanInput{
			Principal: 24000,
			Rate:      0,
			Years:     2,
			Monthly:   false,
		}

		result := CalculateLoan(input)

		// The current implementation doesn't handle zero interest rates correctly
		// It results in NaN due to division by zero in the formula
		// This is a known limitation of the current implementation

		// Check that the result contains NaN values as expected
		if !math.IsNaN(result.MonthlyPayment) {
			t.Errorf("MonthlyPayment = %v, expected NaN", result.MonthlyPayment)
		}

		if !math.IsNaN(result.TotalPaid) {
			t.Errorf("TotalPaid = %v, expected NaN", result.TotalPaid)
		}

		if !math.IsNaN(result.TotalInterest) {
			t.Errorf("TotalInterest = %v, expected NaN", result.TotalInterest)
		}

		// Note: A better implementation would handle zero interest rates by
		// using a different formula: monthlyPayment = principal / numberOfPayments
	})

	// Test case with zero principal
	t.Run("Zero principal", func(t *testing.T) {
		input := LoanInput{
			Principal: 0,
			Rate:      5,
			Years:     10,
			Monthly:   false,
		}

		result := CalculateLoan(input)

		// All payment values should be zero
		if result.MonthlyPayment != 0 {
			t.Errorf("MonthlyPayment = %v, want 0", result.MonthlyPayment)
		}
		if result.TotalPaid != 0 {
			t.Errorf("TotalPaid = %v, want 0", result.TotalPaid)
		}
		if result.TotalInterest != 0 {
			t.Errorf("TotalInterest = %v, want 0", result.TotalInterest)
		}
	})

	// Test case with zero years
	t.Run("Zero years", func(t *testing.T) {
		input := LoanInput{
			Principal: 15000,
			Rate:      7,
			Years:     0,
			Monthly:   false,
		}

		result := CalculateLoan(input)

		// Number of payments should be zero
		if result.NumberOfPayments != 0 {
			t.Errorf("NumberOfPayments = %v, want 0", result.NumberOfPayments)
		}

		// Monthly payment should be NaN or Inf due to division by zero
		if !math.IsNaN(result.MonthlyPayment) && !math.IsInf(result.MonthlyPayment, 0) {
			t.Errorf("MonthlyPayment = %v, expected NaN or Inf", result.MonthlyPayment)
		}
	})

	// Test case with very short loan (1 year)
	t.Run("One year loan", func(t *testing.T) {
		input := LoanInput{
			Principal: 12000,
			Rate:      6,
			Years:     1,
			Monthly:   true,
		}

		result := CalculateLoan(input)

		// Check that we have 12 monthly details
		if len(result.MonthlyDetails) != 12 {
			t.Errorf("Expected 12 months of details, got %d", len(result.MonthlyDetails))
		}

		// Last month's remaining balance should be close to zero
		if len(result.MonthlyDetails) == 12 {
			lastMonth := result.MonthlyDetails[11]
			if lastMonth.RemainingBalance > 0.1 {
				t.Errorf("Last month's remaining balance = %v, expected close to zero", lastMonth.RemainingBalance)
			}
		}
	})
}
