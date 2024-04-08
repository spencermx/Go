package main

import (
	"fmt"
	"math"
)

func main() {
	result := calculate(5, 3)
	fmt.Printf("The result is: %.2f\n", result)
}

func calculate(a, b float64) float64 {
	sum := a + b
	diff := a - b
	prod := a * b
	quot := a / b

	result := sum + diff + prod + quot
	result = math.Sqrt(result)

	return result
}
