package internal

// BudgetInput represents the input parameters for budget allocation
type BudgetInput struct {
	Income        float64
	Housing       float64
	Food          float64
	Transport     float64
	Utilities     float64
	Healthcare    float64
	Debt          float64
	Savings       float64
	Discretionary float64
}

// BudgetCategory represents a single budget category
type BudgetCategory struct {
	Name       string
	Amount     float64
	Percentage float64
}

// BudgetResult represents the output of budget allocation
type BudgetResult struct {
	Income          float64
	Categories      []BudgetCategory
	Total           float64
	TotalPercentage float64
	Warning         string
}

func AllocateBudget(input BudgetInput) BudgetResult {
	totalPercentage := input.Housing + input.Food + input.Transport + input.Utilities +
		input.Healthcare + input.Debt + input.Savings + input.Discretionary

	result := BudgetResult{
		Income:          input.Income,
		TotalPercentage: totalPercentage,
		Categories:      []BudgetCategory{},
	}

	if totalPercentage != 100 {
		result.Warning = "Warning: Your budget percentages total does not equal 100%"
	}

	// Add all budget categories
	result.Categories = []BudgetCategory{
		{
			Name:       "Housing",
			Amount:     input.Income * input.Housing / 100,
			Percentage: input.Housing,
		},
		{
			Name:       "Food",
			Amount:     input.Income * input.Food / 100,
			Percentage: input.Food,
		},
		{
			Name:       "Transportation",
			Amount:     input.Income * input.Transport / 100,
			Percentage: input.Transport,
		},
		{
			Name:       "Utilities",
			Amount:     input.Income * input.Utilities / 100,
			Percentage: input.Utilities,
		},
		{
			Name:       "Healthcare",
			Amount:     input.Income * input.Healthcare / 100,
			Percentage: input.Healthcare,
		},
		{
			Name:       "Debt Repayment",
			Amount:     input.Income * input.Debt / 100,
			Percentage: input.Debt,
		},
		{
			Name:       "Savings",
			Amount:     input.Income * input.Savings / 100,
			Percentage: input.Savings,
		},
		{
			Name:       "Discretionary",
			Amount:     input.Income * input.Discretionary / 100,
			Percentage: input.Discretionary,
		},
	}
	result.Total = input.Income * totalPercentage / 100

	return result
}
