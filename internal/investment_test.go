package internal

import (
	"testing"
)

func TestCalculateInvestment(t *testing.T) {
	tests := []struct {
		name     string
		input    InvestmentInput
		expected InvestmentResult
	}{
		{
			name: "Basic investment calculation",
			input: InvestmentInput{
				Principal:   10000,
				AnnualYield: 7,
				TaxRate:     20,
				Inflation:   2,
				Years:       10,
			},
			expected: InvestmentResult{
				Principal:      10000,
				NetFutureValue: 17865.70, // 10000 * (1.07^10) - tax
				RealValue:      14661.23, // Adjusted for 2% inflation over 10 years
				TaxPaid:        1934.30,  // Tax on profit
				Years:          10,
			},
		},
		{
			name: "Zero yield investment",
			input: InvestmentInput{
				Principal:   5000,
				AnnualYield: 0,
				TaxRate:     15,
				Inflation:   3,
				Years:       5,
			},
			expected: InvestmentResult{
				Principal:      5000,
				NetFutureValue: 5000,    // No growth
				RealValue:      4310.75, // Adjusted for 3% inflation over 5 years
				TaxPaid:        0,       // No profit, no tax
				Years:          5,
			},
		},
		{
			name: "High yield, long term investment",
			input: InvestmentInput{
				Principal:   1000,
				AnnualYield: 12,
				TaxRate:     25,
				Inflation:   2.5,
				Years:       30,
			},
			expected: InvestmentResult{
				Principal:      1000,
				NetFutureValue: 22719.94, // 1000 * (1.12^30) - tax
				RealValue:      10831.57, // Adjusted for 2.5% inflation over 30 years
				TaxPaid:        7239.98,  // Tax on profit
				Years:          30,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateInvestment(tc.input)

			// Check basic fields
			if result.Principal != tc.expected.Principal {
				t.Errorf("Principal = %v, want %v", result.Principal, tc.expected.Principal)
			}
			if result.Years != tc.expected.Years {
				t.Errorf("Years = %v, want %v", result.Years, tc.expected.Years)
			}

			// For floating point values, use approximate comparison
			const tolerance = 0.01 // 1% tolerance

			if !approximatelyEqual(result.NetFutureValue, tc.expected.NetFutureValue, tolerance) {
				t.Errorf("NetFutureValue = %v, want approximately %v", result.NetFutureValue, tc.expected.NetFutureValue)
			}
			if !approximatelyEqual(result.RealValue, tc.expected.RealValue, tolerance) {
				t.Errorf("RealValue = %v, want approximately %v", result.RealValue, tc.expected.RealValue)
			}
			if !approximatelyEqual(result.TaxPaid, tc.expected.TaxPaid, tolerance) {
				t.Errorf("TaxPaid = %v, want approximately %v", result.TaxPaid, tc.expected.TaxPaid)
			}
		})
	}
}

// Using approximatelyEqual function from retirement_test.go

// TestInvestmentEdgeCases tests edge cases for the investment calculator
func TestInvestmentEdgeCases(t *testing.T) {
	// Test case with zero principal
	t.Run("Zero principal", func(t *testing.T) {
		input := InvestmentInput{
			Principal:   0,
			AnnualYield: 8,
			TaxRate:     20,
			Inflation:   2,
			Years:       10,
		}

		result := CalculateInvestment(input)

		if result.NetFutureValue != 0 {
			t.Errorf("NetFutureValue = %v, want 0", result.NetFutureValue)
		}
		if result.RealValue != 0 {
			t.Errorf("RealValue = %v, want 0", result.RealValue)
		}
		if result.TaxPaid != 0 {
			t.Errorf("TaxPaid = %v, want 0", result.TaxPaid)
		}
	})

	// Test case with negative yield (market downturn)
	t.Run("Negative yield", func(t *testing.T) {
		input := InvestmentInput{
			Principal:   10000,
			AnnualYield: -5,
			TaxRate:     20,
			Inflation:   2,
			Years:       5,
		}

		result := CalculateInvestment(input)

		// With negative yield, future value should be less than principal
		if result.NetFutureValue >= input.Principal {
			t.Errorf("With negative yield, NetFutureValue (%v) should be less than Principal (%v)",
				result.NetFutureValue, input.Principal)
		}

		// Real value should be even less due to inflation
		if result.RealValue >= result.NetFutureValue {
			t.Errorf("RealValue (%v) should be less than NetFutureValue (%v) due to inflation",
				result.RealValue, result.NetFutureValue)
		}
	})

	// Test case with zero years
	t.Run("Zero years", func(t *testing.T) {
		input := InvestmentInput{
			Principal:   5000,
			AnnualYield: 6,
			TaxRate:     15,
			Inflation:   3,
			Years:       0,
		}

		result := CalculateInvestment(input)

		// With zero years, future value should equal principal
		if result.NetFutureValue != input.Principal {
			t.Errorf("NetFutureValue = %v, want %v", result.NetFutureValue, input.Principal)
		}
		if result.RealValue != input.Principal {
			t.Errorf("RealValue = %v, want %v", result.RealValue, input.Principal)
		}
		if result.TaxPaid != 0 {
			t.Errorf("TaxPaid = %v, want 0", result.TaxPaid)
		}
	})
}
