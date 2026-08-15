// Package calc provides basic arithmetic operations for the calculator CLI.
//
// All operations accept float64 operands and return a float64 result.
// Division by zero returns an error rather than panicking.
package calc

import "fmt"

// Add returns the sum of two operands.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns the difference of a minus b.
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns the product of two operands.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns the quotient of a divided by b.
// It returns an error if b is zero to prevent a panic.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero is not allowed")
	}
	return a / b, nil
}
