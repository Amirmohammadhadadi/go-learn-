package __Functions

import "fmt"

func add(x, y int) int {
	return x + y
}
func mul(x, y int) int {
	return x * y
}
func aggregate(a, b, c int, arithmetic func(int, int) int) int {
	return arithmetic(arithmetic(a, b), c)
}
func _main() {
	fmt.Println(aggregate(2, 3, 4, add))
	fmt.Println(aggregate(2, 3, 4, mul))
}
