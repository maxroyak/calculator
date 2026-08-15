// Calculator is a simple CLI tool that performs basic arithmetic operations:
// add, subtract, multiply, and divide.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/maxroyak/calculator/internal/calc"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return printUsage()
	}

	op := args[0]
	if op == "-h" || op == "--help" {
		return printUsage()
	}

	if len(args) < 3 {
		return fmt.Errorf("operation %q requires two numeric operands", op)
	}

	a, err := parseFloat(args[1])
	if err != nil {
		return fmt.Errorf("first operand: %w", err)
	}

	b, err := parseFloat(args[2])
	if err != nil {
		return fmt.Errorf("second operand: %w", err)
	}

	result, err := calculate(op, a, b)
	if err != nil {
		return err
	}

	fmt.Printf("%g\n", result)
	return nil
}

func calculate(op string, a, b float64) (float64, error) {
	switch op {
	case "add":
		return calc.Add(a, b), nil
	case "subtract":
		return calc.Subtract(a, b), nil
	case "multiply":
		return calc.Multiply(a, b), nil
	case "divide":
		return calc.Divide(a, b)
	default:
		return 0, fmt.Errorf("unknown operation %q: valid operations are add, subtract, multiply, divide", op)
	}
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%g", &f)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", s)
	}
	return f, nil
}

func printUsage() error {
	fmt.Println(`calculator - a simple CLI calculator

Usage:
  calculator <operation> <operand1> <operand2>

Operations:
  add       Add two numbers
  subtract  Subtract second number from first
  multiply  Multiply two numbers
  divide    Divide first number by second

Examples:
  calculator add 5 3        # Output: 8
  calculator subtract 10 4  # Output: 6
  calculator multiply 6 7   # Output: 42
  calculator divide 10 2    # Output: 5`)
	return errors.New("no operation provided")
}
