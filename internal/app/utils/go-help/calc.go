package help

import (
	"fmt"
	"strconv"
)

// Add func
func Add(num1, num2 float64) float64 {
	return num1 + num2
}

// Subtract func
func Subtract(num1, num2 float64) float64 {
	return num1 - num2
}

// Multiply func
func Multiply(num1, num2 float64) float64 {
	return ToFixed(num1 * num2)
}

// Divide func
func Divide(num1, num2 float64) float64 {
	if num2 == float64(0) {
		return 0
	}
	return ToFixed(num1 / num2)
}

// ToFixed func
func ToFixed(v float64) float64 {
	// FLOAT TO STRING FORMATING DIGIT DECIMAL
	res := fmt.Sprintf("%.4f", v)
	// STRING TO FLOAT64
	v, _ = strconv.ParseFloat(res, 64)

	return v
}

// Division func
func Division(term, idx int, val float64) float64 {
	// MAIN VARIABLE
	var result []float64
	var totalVal float64

	// GET THE RESULT INDEXING
	for i := 1; i <= term; i++ {
		var calc float64
		if i == term {
			calc = val - totalVal
		} else {
			calc = val / float64(term)
			ToFixed(calc)
		}
		totalVal += calc

		result = append(result, calc)
	}

	// SEND THE RESULT INDEXING
	return result[idx]
}
