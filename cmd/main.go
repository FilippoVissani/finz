package main

import (
	"finz/internal"
	"flag"
	"fmt"
	"os"
	"strings"
)

func handleInvest(args []string) {
	investCmd := flag.NewFlagSet("invest", flag.ExitOnError)

	var (
		principal   float64
		annualYield float64
		taxRate     float64
		inflation   float64
		years       int
	)

	investCmd.Float64Var(&principal, "initial", 10000, "Initial investment amount")
	investCmd.Float64Var(&annualYield, "yield", 7.0, "Annual yield in percent (e.g., 7)")
	investCmd.Float64Var(&taxRate, "tax", 26.0, "Tax rate on gains in percent")
	investCmd.Float64Var(&inflation, "inflation", 2.0, "Annual inflation rate in percent")
	investCmd.IntVar(&years, "years", 10, "Investment duration in years")

	if err := investCmd.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if investCmd.Parsed() {
		if investCmd.NArg() > 0 {
			fmt.Printf("Unknown arguments: %s\n", strings.Join(investCmd.Args(), " "))
			investCmd.PrintDefaults()
			os.Exit(1)
		}
	}

	input := internal.InvestmentInput{
		Principal:   principal,
		AnnualYield: annualYield,
		TaxRate:     taxRate,
		Inflation:   inflation,
		Years:       years,
	}

	result := internal.CalculateInvestment(input)

	fmt.Printf("Initial amount:        €%.2f\n", result.Principal)
	fmt.Printf("Nominal final value:   €%.2f\n", result.NetFutureValue)
	fmt.Printf("Real final value:      €%.2f\n", result.RealValue)
	fmt.Printf("Total tax paid:        €%.2f\n", result.TaxPaid)
	fmt.Printf("Total years:           %d\n", result.Years)
}

func handleLoan(args []string) {
	loanCmd := flag.NewFlagSet("loan", flag.ExitOnError)

	var (
		principal float64
		rate      float64
		years     int
		monthly   bool
	)

	loanCmd.Float64Var(&principal, "amount", 100000, "Loan amount")
	loanCmd.Float64Var(&rate, "rate", 4.5, "Annual interest rate in percent")
	loanCmd.IntVar(&years, "years", 30, "Loan term in years")
	loanCmd.BoolVar(&monthly, "monthly", true, "Show monthly payment breakdown")

	if err := loanCmd.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if loanCmd.Parsed() {
		if loanCmd.NArg() > 0 {
			fmt.Printf("Unknown arguments: %s\n", strings.Join(loanCmd.Args(), " "))
			loanCmd.PrintDefaults()
			os.Exit(1)
		}
	}

	input := internal.LoanInput{
		Principal: principal,
		Rate:      rate,
		Years:     years,
		Monthly:   monthly,
	}

	result := internal.CalculateLoan(input)

	fmt.Printf("Loan amount:           €%.2f\n", result.Principal)
	fmt.Printf("Monthly payment:       €%.2f\n", result.MonthlyPayment)
	fmt.Printf("Total paid:            €%.2f\n", result.TotalPaid)
	fmt.Printf("Total interest:        €%.2f\n", result.TotalInterest)
	fmt.Printf("Loan term:             %d years (%d payments)\n", result.Years, result.NumberOfPayments)

	if monthly && len(result.MonthlyDetails) > 0 {
		fmt.Println("\nMonthly Payment Breakdown:")
		fmt.Println("Month\tPayment\t\tPrincipal\tInterest\tRemaining")

		for _, detail := range result.MonthlyDetails {
			fmt.Printf("%d\t€%.2f\t\t€%.2f\t\t€%.2f\t\t€%.2f\n",
				detail.Month, detail.Payment, detail.PrincipalPayment, detail.InterestPayment, detail.RemainingBalance)
		}
		fmt.Println("...")
	}
}

func handleSavings(args []string) {
	savingsCmd := flag.NewFlagSet("savings", flag.ExitOnError)

	var (
		initial        float64
		monthlyDeposit float64
		annualYield    float64
		inflation      float64
		years          int
	)

	savingsCmd.Float64Var(&initial, "initial", 1000, "Initial deposit amount")
	savingsCmd.Float64Var(&monthlyDeposit, "monthly", 100, "Monthly deposit amount")
	savingsCmd.Float64Var(&annualYield, "yield", 3.0, "Annual yield in percent")
	savingsCmd.Float64Var(&inflation, "inflation", 2.0, "Annual inflation rate in percent")
	savingsCmd.IntVar(&years, "years", 10, "Savings duration in years")

	if err := savingsCmd.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if savingsCmd.Parsed() {
		if savingsCmd.NArg() > 0 {
			fmt.Printf("Unknown arguments: %s\n", strings.Join(savingsCmd.Args(), " "))
			savingsCmd.PrintDefaults()
			os.Exit(1)
		}
	}

	input := internal.SavingsInput{
		Initial:        initial,
		MonthlyDeposit: monthlyDeposit,
		AnnualYield:    annualYield,
		Inflation:      inflation,
		Years:          years,
	}

	result := internal.CalculateSavings(input)

	fmt.Printf("Initial deposit:       €%.2f\n", result.Initial)
	fmt.Printf("Monthly deposit:       €%.2f\n", result.MonthlyDeposit)
	fmt.Printf("Nominal final balance: €%.2f\n", result.FutureValue)
	fmt.Printf("Real final balance:    €%.2f\n", result.RealFutureValue)
	fmt.Printf("Total deposits:        €%.2f\n", result.TotalDeposits)
	fmt.Printf("Interest earned:       €%.2f\n", result.InterestEarned)
	fmt.Printf("Savings period:        %d years (%d months)\n", result.Years, result.NumberOfMonths)
}

func handleRetirement(args []string) {
	retireCmd := flag.NewFlagSet("retirement", flag.ExitOnError)

	var (
		currentAge          int
		retirementAge       int
		currentSavings      float64
		monthlyContribution float64
		withdrawalRate      float64
		annualYield         float64
		inflation           float64
	)

	retireCmd.IntVar(&currentAge, "age", 30, "Current age")
	retireCmd.IntVar(&retirementAge, "retire-age", 65, "Retirement age")
	retireCmd.Float64Var(&currentSavings, "savings", 50000, "Current retirement savings")
	retireCmd.Float64Var(&monthlyContribution, "monthly", 500, "Monthly contribution")
	retireCmd.Float64Var(&withdrawalRate, "withdrawal", 4.0, "Annual withdrawal rate in percent")
	retireCmd.Float64Var(&annualYield, "yield", 7.0, "Annual investment yield in percent")
	retireCmd.Float64Var(&inflation, "inflation", 2.0, "Annual inflation rate in percent")

	if err := retireCmd.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if retireCmd.Parsed() {
		if retireCmd.NArg() > 0 {
			fmt.Printf("Unknown arguments: %s\n", strings.Join(retireCmd.Args(), " "))
			retireCmd.PrintDefaults()
			os.Exit(1)
		}
	}

	input := internal.RetirementInput{
		CurrentAge:          currentAge,
		RetirementAge:       retirementAge,
		CurrentSavings:      currentSavings,
		MonthlyContribution: monthlyContribution,
		WithdrawalRate:      withdrawalRate,
		AnnualYield:         annualYield,
		Inflation:           inflation,
	}

	result := internal.CalculateRetirement(input)

	fmt.Printf("Current age:           %d\n", result.CurrentAge)
	fmt.Printf("Retirement age:        %d\n", result.RetirementAge)
	fmt.Printf("Years to retirement:   %d\n", result.YearsToRetirement)
	fmt.Printf("Retirement savings:    €%.2f\n", result.RetirementSavings)
	fmt.Printf("Annual withdrawal:     €%.2f\n", result.AnnualWithdrawal)
	fmt.Printf("Monthly withdrawal:    €%.2f\n", result.MonthlyWithdrawal)
	fmt.Printf("Inflation-adjusted monthly withdrawal: €%.2f\n", result.RealMonthlyWithdrawal)
}

func handleCurrency(args []string) {
	currencyCmd := flag.NewFlagSet("currency", flag.ExitOnError)

	var (
		amount float64
		from   string
		to     string
	)

	currencyCmd.Float64Var(&amount, "amount", 100, "Amount to convert")
	currencyCmd.StringVar(&from, "from", "EUR", "Source currency code (e.g., EUR, USD)")
	currencyCmd.StringVar(&to, "to", "USD", "Target currency code (e.g., EUR, USD)")

	if err := currencyCmd.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if currencyCmd.Parsed() {
		if currencyCmd.NArg() > 0 {
			fmt.Printf("Unknown arguments: %s\n", strings.Join(currencyCmd.Args(), " "))
			currencyCmd.PrintDefaults()
			os.Exit(1)
		}
	}

	input := internal.CurrencyInput{
		Amount: amount,
		From:   from,
		To:     to,
	}

	result := internal.ConvertCurrency(input)

	if result.Error != nil {
		fmt.Println(result.Error)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", result.Amount, result.From, result.Converted, result.To)
	fmt.Printf("Exchange rate: 1 %s = %.4f %s\n", result.From, result.ExchangeRate, result.To)
}

func handleBudget(args []string) {
	budgetCmd := flag.NewFlagSet("budget", flag.ExitOnError)

	var (
		income        float64
		housing       float64
		food          float64
		transport     float64
		utilities     float64
		healthcare    float64
		debt          float64
		savings       float64
		discretionary float64
	)

	budgetCmd.Float64Var(&income, "income", 3000, "Monthly income")
	budgetCmd.Float64Var(&housing, "housing", 30, "Housing percentage")
	budgetCmd.Float64Var(&food, "food", 15, "Food percentage")
	budgetCmd.Float64Var(&transport, "transport", 10, "Transportation percentage")
	budgetCmd.Float64Var(&utilities, "utilities", 5, "Utilities percentage")
	budgetCmd.Float64Var(&healthcare, "healthcare", 5, "Healthcare percentage")
	budgetCmd.Float64Var(&debt, "debt", 10, "Debt repayment percentage")
	budgetCmd.Float64Var(&savings, "savings", 15, "Savings percentage")
	budgetCmd.Float64Var(&discretionary, "discretionary", 10, "Discretionary spending percentage")

	if err := budgetCmd.Parse(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if budgetCmd.Parsed() {
		if budgetCmd.NArg() > 0 {
			fmt.Printf("Unknown arguments: %s\n", strings.Join(budgetCmd.Args(), " "))
			budgetCmd.PrintDefaults()
			os.Exit(1)
		}
	}

	input := internal.BudgetInput{
		Income:        income,
		Housing:       housing,
		Food:          food,
		Transport:     transport,
		Utilities:     utilities,
		Healthcare:    healthcare,
		Debt:          debt,
		Savings:       savings,
		Discretionary: discretionary,
	}

	result := internal.AllocateBudget(input)

	if result.Warning != "" {
		fmt.Printf("%s\n\n", result.Warning)
	}

	fmt.Printf("Monthly Income: €%.2f\n\n", result.Income)
	fmt.Printf("Budget Allocation:\n")

	for _, category := range result.Categories {
		fmt.Printf("%-15s €%.2f (%.1f%%)\n", category.Name+":", category.Amount, category.Percentage)
	}

	fmt.Printf("\nTotal:         €%.2f (%.1f%%)\n", result.Total, result.TotalPercentage)
}

func handleHelp() {
	internal.PrintUsage()
}

func handleUnknownCommand(command string) {
	fmt.Printf("Unknown command: %s\n", command)
	internal.PrintUsage()
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		internal.PrintUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "invest":
		handleInvest(args)
	case "loan":
		handleLoan(args)
	case "savings":
		handleSavings(args)
	case "retirement":
		handleRetirement(args)
	case "currency":
		handleCurrency(args)
	case "budget":
		handleBudget(args)
	case "help":
		handleHelp()
	default:
		handleUnknownCommand(command)
	}
}
