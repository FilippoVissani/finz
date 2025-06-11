package internal

import "fmt"

func PrintUsage() {
	fmt.Println("Finz - Financial Calculator CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  finz <command> [options]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  invest      - Calculate investment growth with taxes and inflation")
	fmt.Println("  loan        - Calculate loan or mortgage payments")
	fmt.Println("  savings     - Calculate savings with regular deposits")
	fmt.Println("  retirement  - Calculate retirement savings and withdrawals")
	fmt.Println("  currency    - Convert between currencies")
	fmt.Println("  budget      - Allocate budget based on percentages")
	fmt.Println("  help        - Show this help message")
	fmt.Println("\nRun 'finz <command> --help' for more information on a command.")
}
