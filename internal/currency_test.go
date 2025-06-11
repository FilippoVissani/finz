package internal

import (
	"testing"
)

func TestConvertCurrency(t *testing.T) {
	tests := []struct {
		name     string
		input    CurrencyInput
		expected CurrencyResult
	}{
		{
			name: "EUR to USD conversion",
			input: CurrencyInput{
				Amount: 100,
				From:   "EUR",
				To:     "USD",
			},
			expected: CurrencyResult{
				Amount:       100,
				From:         "EUR",
				To:           "USD",
				Converted:    109,
				ExchangeRate: 1.09,
				Error:        nil,
			},
		},
		{
			name: "USD to JPY conversion",
			input: CurrencyInput{
				Amount: 50,
				From:   "USD",
				To:     "JPY",
			},
			expected: CurrencyResult{
				Amount:       50,
				From:         "USD",
				To:           "JPY",
				Converted:    7350,
				ExchangeRate: 147.0,
				Error:        nil,
			},
		},
		{
			name: "GBP to EUR conversion",
			input: CurrencyInput{
				Amount: 75,
				From:   "GBP",
				To:     "EUR",
			},
			expected: CurrencyResult{
				Amount:       75,
				From:         "GBP",
				To:           "EUR",
				Converted:    88.5,
				ExchangeRate: 1.18,
				Error:        nil,
			},
		},
		{
			name: "JPY to GBP conversion",
			input: CurrencyInput{
				Amount: 10000,
				From:   "JPY",
				To:     "GBP",
			},
			expected: CurrencyResult{
				Amount:       10000,
				From:         "JPY",
				To:           "GBP",
				Converted:    53,
				ExchangeRate: 0.0053,
				Error:        nil,
			},
		},
		{
			name: "Same currency conversion",
			input: CurrencyInput{
				Amount: 200,
				From:   "USD",
				To:     "USD",
			},
			expected: CurrencyResult{
				Amount:       200,
				From:         "USD",
				To:           "USD",
				Converted:    200,
				ExchangeRate: 1.0,
				Error:        nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertCurrency(tc.input)

			// Check basic fields
			if result.Amount != tc.expected.Amount {
				t.Errorf("Amount = %v, want %v", result.Amount, tc.expected.Amount)
			}
			if result.From != tc.expected.From {
				t.Errorf("From = %v, want %v", result.From, tc.expected.From)
			}
			if result.To != tc.expected.To {
				t.Errorf("To = %v, want %v", result.To, tc.expected.To)
			}

			// For floating point values, use approximate comparison
			const tolerance = 0.001 // 0.1% tolerance

			if !approximatelyEqual(result.Converted, tc.expected.Converted, tolerance) {
				t.Errorf("Converted = %v, want approximately %v", result.Converted, tc.expected.Converted)
			}
			if !approximatelyEqual(result.ExchangeRate, tc.expected.ExchangeRate, tolerance) {
				t.Errorf("ExchangeRate = %v, want approximately %v", result.ExchangeRate, tc.expected.ExchangeRate)
			}

			// Check error
			if (result.Error != nil) != (tc.expected.Error != nil) {
				t.Errorf("Error = %v, want %v", result.Error, tc.expected.Error)
			}
		})
	}
}

// TestCurrencyEdgeCases tests edge cases for the currency converter
func TestCurrencyEdgeCases(t *testing.T) {
	// Test case with unsupported source currency
	t.Run("Unsupported source currency", func(t *testing.T) {
		input := CurrencyInput{
			Amount: 100,
			From:   "XYZ", // Invalid currency
			To:     "USD",
		}

		result := ConvertCurrency(input)

		if result.Error == nil {
			t.Errorf("Expected an error for unsupported source currency, got nil")
		} else if result.Error.Error() != "unsupported source currency: XYZ" {
			t.Errorf("Error = %v, want 'unsupported source currency: XYZ'", result.Error)
		}

		// Converted amount and exchange rate should be zero
		if result.Converted != 0 {
			t.Errorf("Converted = %v, want 0", result.Converted)
		}
		if result.ExchangeRate != 0 {
			t.Errorf("ExchangeRate = %v, want 0", result.ExchangeRate)
		}
	})

	// Test case with unsupported target currency
	t.Run("Unsupported target currency", func(t *testing.T) {
		input := CurrencyInput{
			Amount: 100,
			From:   "USD",
			To:     "XYZ", // Invalid currency
		}

		result := ConvertCurrency(input)

		if result.Error == nil {
			t.Errorf("Expected an error for unsupported target currency, got nil")
		} else if result.Error.Error() != "unsupported target currency: XYZ" {
			t.Errorf("Error = %v, want 'unsupported target currency: XYZ'", result.Error)
		}

		// Converted amount and exchange rate should be zero
		if result.Converted != 0 {
			t.Errorf("Converted = %v, want 0", result.Converted)
		}
		if result.ExchangeRate != 0 {
			t.Errorf("ExchangeRate = %v, want 0", result.ExchangeRate)
		}
	})

	// Test case with zero amount
	t.Run("Zero amount", func(t *testing.T) {
		input := CurrencyInput{
			Amount: 0,
			From:   "EUR",
			To:     "USD",
		}

		result := ConvertCurrency(input)

		// No error should occur
		if result.Error != nil {
			t.Errorf("Error = %v, want nil", result.Error)
		}

		// Converted amount should be zero
		if result.Converted != 0 {
			t.Errorf("Converted = %v, want 0", result.Converted)
		}

		// Exchange rate should be correct
		if !approximatelyEqual(result.ExchangeRate, 1.09, 0.001) {
			t.Errorf("ExchangeRate = %v, want 1.09", result.ExchangeRate)
		}
	})

	// Test case with lowercase currency codes
	t.Run("Lowercase currency codes", func(t *testing.T) {
		input := CurrencyInput{
			Amount: 100,
			From:   "eur",
			To:     "usd",
		}

		result := ConvertCurrency(input)

		// No error should occur
		if result.Error != nil {
			t.Errorf("Error = %v, want nil", result.Error)
		}

		// Currency codes should be uppercase in the result
		if result.From != "EUR" {
			t.Errorf("From = %v, want EUR", result.From)
		}
		if result.To != "USD" {
			t.Errorf("To = %v, want USD", result.To)
		}

		// Conversion should work correctly
		if !approximatelyEqual(result.Converted, 109, 0.001) {
			t.Errorf("Converted = %v, want 109", result.Converted)
		}
	})
}
