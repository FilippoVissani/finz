# Finz

A comprehensive financial calculator CLI tool written in Go.

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/finz.git
cd finz

# Build the project
make build

# The binary will be available in .build/finz
```

## Usage

```bash
finz <command> [options]
```

### Available Commands

- `invest` - Calculate investment growth with taxes and inflation
- `loan` - Calculate loan or mortgage payments
- `savings` - Calculate savings with regular deposits
- `retirement` - Calculate retirement savings and withdrawals
- `currency` - Convert between currencies
- `budget` - Allocate budget based on percentages
- `help` - Show help message

## Examples

### Investment Calculator

Calculate the future value of an investment with taxes and inflation:

```bash
finz invest --initial 10000 --yield 7 --tax 26 --inflation 2 --years 10
```

### Loan Calculator

Calculate mortgage or loan payments:

```bash
finz loan --amount 250000 --rate 3.5 --years 30 --monthly
```

### Savings Calculator

Calculate savings growth with regular deposits:

```bash
finz savings --initial 1000 --monthly 200 --yield 3 --years 5
```

### Retirement Calculator

Plan for retirement:

```bash
finz retirement --age 35 --retire-age 65 --savings 50000 --monthly 500 --yield 7 --inflation 2
```

### Currency Converter

Convert between currencies:

```bash
finz currency --amount 100 --from EUR --to USD
```

### Budget Allocation

Allocate your monthly budget:

```bash
finz budget --income 3000 --housing 30 --food 15 --transport 10 --utilities 5 --healthcare 5 --debt 10 --savings 15 --discretionary 10
```
