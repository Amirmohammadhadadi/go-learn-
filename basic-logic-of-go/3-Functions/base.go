package __Functions


import (
	"fmt"
)

func main() {
	x, _ := getPoint()
}

func sub(x int, y int) int {
	return x - y
}

//	func addToDatabase(hp, damage int) {
//		// ...
//	}
//
//	func addToDatabase(hp, damage int, name string) {
//		// ?
//	}
//
//	func addToDatabase(hp, damage int, name string, level int) {
//		// ?
//	}
func multi(callback func(int, int) int, number int) int {
	return callback(number, 2)
}

//f func(func(int,int) int, int) int


func main() {
   x := 5
   increment(x)

   fmt.Println(x)
   // still prints 5,
   // because the increment function
   // received a copy of x
}

func increment(x int) {
   x++
}


func getPoint() (x int, y int) {
	return 3, 4
}

x, _ := getPoint()
//ignore y value


func getCoords() (x, y int) {
	// x and y are initialized with zero values

	return // automatically returns x and y
}
//Is the same as:
func getCoords() (int, int) {
	var x int
	var y int
	return x, y
}


