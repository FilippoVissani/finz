package internal

import (
	"testing"
)

func TestAllocateBudget(t *testing.T) {
	tests := []struct {
		name     string
		input    BudgetInput
		expected BudgetResult
	}{
		{
			name: "Basic budget allocation with 100% total",
			input: BudgetInput{
				Income:        5000,
				Housing:       30,
				Food:          15,
				Transport:     10,
				Utilities:     5,
				Healthcare:    10,
				Debt:          10,
				Savings:       15,
				Discretionary: 5,
			},
			expected: BudgetResult{
				Income:          5000,
				TotalPercentage: 100,
				Total:           5000,
				Warning:         "",
			},
		},
		{
			name: "Budget allocation with less than 100% total",
			input: BudgetInput{
				Income:        3000,
				Housing:       25,
				Food:          15,
				Transport:     10,
				Utilities:     5,
				Healthcare:    5,
				Debt:          5,
				Savings:       10,
				Discretionary: 5,
			},
			expected: BudgetResult{
				Income:          3000,
				TotalPercentage: 80,
				Total:           2400,
				Warning:         "Warning: Your budget percentages total does not equal 100%",
			},
		},
		{
			name: "Budget allocation with more than 100% total",
			input: BudgetInput{
				Income:        4000,
				Housing:       35,
				Food:          20,
				Transport:     15,
				Utilities:     10,
				Healthcare:    10,
				Debt:          15,
				Savings:       10,
				Discretionary: 5,
			},
			expected: BudgetResult{
				Income:          4000,
				TotalPercentage: 120,
				Total:           4800,
				Warning:         "Warning: Your budget percentages total does not equal 100%",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := AllocateBudget(tc.input)

			// Check basic fields
			if result.Income != tc.expected.Income {
				t.Errorf("Income = %v, want %v", result.Income, tc.expected.Income)
			}

			if !approximatelyEqual(result.TotalPercentage, tc.expected.TotalPercentage, 0.001) {
				t.Errorf("TotalPercentage = %v, want %v", result.TotalPercentage, tc.expected.TotalPercentage)
			}

			if !approximatelyEqual(result.Total, tc.expected.Total, 0.001) {
				t.Errorf("Total = %v, want %v", result.Total, tc.expected.Total)
			}

			if result.Warning != tc.expected.Warning {
				t.Errorf("Warning = %v, want %v", result.Warning, tc.expected.Warning)
			}

			// Check that we have the expected number of categories
			if len(result.Categories) != 8 {
				t.Errorf("Expected 8 budget categories, got %d", len(result.Categories))
			}

			// Check individual category calculations
			checkCategory(t, result, "Housing", tc.input.Income*tc.input.Housing/100, tc.input.Housing)
			checkCategory(t, result, "Food", tc.input.Income*tc.input.Food/100, tc.input.Food)
			checkCategory(t, result, "Transportation", tc.input.Income*tc.input.Transport/100, tc.input.Transport)
			checkCategory(t, result, "Utilities", tc.input.Income*tc.input.Utilities/100, tc.input.Utilities)
			checkCategory(t, result, "Healthcare", tc.input.Income*tc.input.Healthcare/100, tc.input.Healthcare)
			checkCategory(t, result, "Debt Repayment", tc.input.Income*tc.input.Debt/100, tc.input.Debt)
			checkCategory(t, result, "Savings", tc.input.Income*tc.input.Savings/100, tc.input.Savings)
			checkCategory(t, result, "Discretionary", tc.input.Income*tc.input.Discretionary/100, tc.input.Discretionary)
		})
	}
}

// Helper function to check a specific budget category
func checkCategory(t *testing.T, result BudgetResult, name string, expectedAmount, expectedPercentage float64) {
	t.Helper()

	for _, category := range result.Categories {
		if category.Name == name {
			if !approximatelyEqual(category.Amount, expectedAmount, 0.001) {
				t.Errorf("%s amount = %v, want %v", name, category.Amount, expectedAmount)
			}

			if !approximatelyEqual(category.Percentage, expectedPercentage, 0.001) {
				t.Errorf("%s percentage = %v, want %v", name, category.Percentage, expectedPercentage)
			}

			return
		}
	}

	t.Errorf("Category %s not found", name)
}

// TestBudgetEdgeCases tests edge cases for the budget calculator
func TestBudgetEdgeCases(t *testing.T) {
	// Test case with zero income
	t.Run("Zero income", func(t *testing.T) {
		input := BudgetInput{
			Income:        0,
			Housing:       30,
			Food:          15,
			Transport:     10,
			Utilities:     5,
			Healthcare:    10,
			Debt:          10,
			Savings:       15,
			Discretionary: 5,
		}

		result := AllocateBudget(input)

		if result.Income != 0 {
			t.Errorf("Income = %v, want 0", result.Income)
		}

		if result.Total != 0 {
			t.Errorf("Total = %v, want 0", result.Total)
		}

		// Check that all category amounts are zero
		for _, category := range result.Categories {
			if category.Amount != 0 {
				t.Errorf("%s amount = %v, want 0", category.Name, category.Amount)
			}
		}
	})

	// Test case with zero percentages
	t.Run("Zero percentages", func(t *testing.T) {
		input := BudgetInput{
			Income:        3000,
			Housing:       0,
			Food:          0,
			Transport:     0,
			Utilities:     0,
			Healthcare:    0,
			Debt:          0,
			Savings:       0,
			Discretionary: 0,
		}

		result := AllocateBudget(input)

		if result.TotalPercentage != 0 {
			t.Errorf("TotalPercentage = %v, want 0", result.TotalPercentage)
		}

		if result.Total != 0 {
			t.Errorf("Total = %v, want 0", result.Total)
		}

		if result.Warning == "" {
			t.Errorf("Expected a warning message for zero percentages")
		}

		// Check that all category amounts are zero
		for _, category := range result.Categories {
			if category.Amount != 0 {
				t.Errorf("%s amount = %v, want 0", category.Name, category.Amount)
			}

			if category.Percentage != 0 {
				t.Errorf("%s percentage = %v, want 0", category.Name, category.Percentage)
			}
		}
	})

	// Test case with negative income
	t.Run("Negative income", func(t *testing.T) {
		input := BudgetInput{
			Income:        -2000,
			Housing:       30,
			Food:          15,
			Transport:     10,
			Utilities:     5,
			Healthcare:    10,
			Debt:          10,
			Savings:       15,
			Discretionary: 5,
		}

		result := AllocateBudget(input)

		// With negative income, all category amounts should be negative
		for _, category := range result.Categories {
			if category.Amount >= 0 {
				t.Errorf("%s amount = %v, should be negative", category.Name, category.Amount)
			}
		}

		// Total should be negative
		if result.Total >= 0 {
			t.Errorf("Total = %v, should be negative", result.Total)
		}
	})
}
