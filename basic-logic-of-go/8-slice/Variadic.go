package main

import "fmt"

func sum(nums ...int) int {
	// nums is just a slice
	result := 0
	for i := 0; i < len(nums); i++ {
		num := nums[i]
		result += num
	}
	return result
}

func variadic() {
	total := sum(1, 2, 3, 4, 5)
	fmt.Printf("total: %v\n", total)
}

// Spread Operator
func printStrings(strings ...string) []string {
	return strings
}
func printStrings_() {
	names := []string{"a", "b", "c", "d"}
	printStrings(names...)
}
