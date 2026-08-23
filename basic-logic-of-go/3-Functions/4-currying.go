package __Functions

import "fmt"

//Function currying is a concept from functional programming and involves partial application of functions. It allows a function with multiple arguments to be transformed into a sequence of functions, each taking a single argument.
//Let's simulate this behavior. For example:

func main() {
	squareFunc := selfMath(multiply)
	doubleFunc := selfMath(_add)

	fmt.Println(squareFunc(5))
	// prints 25

	fmt.Println(doubleFunc(5))
	// prints 10
}

func multiply(x, y int) int {
	return x * y
}

func _add(x, y int) int {
	return x + y
}

func selfMath(mathFunc func(int, int) int) func(int) int {
	return func(x int) int {
		return mathFunc(x, x)
	}
}

/*
In the example above, the selfMath function takes in a function as its parameter and returns a function that itself returns the value of running that input function on its parameter.
*/
