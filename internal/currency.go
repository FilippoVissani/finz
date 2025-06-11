package internal

import (
	"errors"
	"strings"
)

// CurrencyInput represents the input parameters for currency conversion
type CurrencyInput struct {
	Amount float64
	From   string
	To     string
}

// CurrencyResult represents the output of currency conversion
type CurrencyResult struct {
	Amount       float64
	From         string
	To           string
	Converted    float64
	ExchangeRate float64
	Error        error
}

func ConvertCurrency(input CurrencyInput) CurrencyResult {
	// Simple hardcoded exchange rates (in a real app, these would come from an API)
	rates := map[string]map[string]float64{
		"EUR": {"USD": 1.09, "GBP": 0.85, "JPY": 160.0},
		"USD": {"EUR": 0.92, "GBP": 0.78, "JPY": 147.0},
		"GBP": {"EUR": 1.18, "USD": 1.28, "JPY": 188.0},
		"JPY": {"EUR": 0.00625, "USD": 0.0068, "GBP": 0.0053},
	}

	from := strings.ToUpper(input.From)
	to := strings.ToUpper(input.To)

	result := CurrencyResult{
		Amount: input.Amount,
		From:   from,
		To:     to,
	}

	if from == to {
		result.Converted = input.Amount
		result.ExchangeRate = 1.0
		return result
	}

	if _, exists := rates[from]; !exists {
		result.Error = errors.New("unsupported source currency: " + from)
		return result
	}

	if rate, exists := rates[from][to]; exists {
		result.Converted = input.Amount * rate
		result.ExchangeRate = rate
		return result
	} else {
		result.Error = errors.New("unsupported target currency: " + to)
		return result
	}
}
