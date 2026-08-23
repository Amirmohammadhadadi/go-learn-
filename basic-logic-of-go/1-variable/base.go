package __variable

import "fmt"

/*
All data 1-variable
Signed integers (no decimal)
int  int8  int16  int32  int64

Unsigned integers (non-negative numbers/no decimal)
uint uint8 uint16 uint32 uint64 uintptr

Signed decimal numbers
float32 float64

Complex numbers (a complex number has a real and imaginary part)
complex64 complex128




bool

string

int int8 int16 int32 int64
uint uint8 uint16 uint32 uint64 uintptr

byte // alias for uint8

rune // alias for uint32
	//represent a Unicode code print

float32 float64

complex64 complex128

*/

func _main() {
	var number int
	var pi float64 = 3.14159
	var hasIt bool = false
	var username string
	var empty string = ""
	empty1 := "Hello World"
	var isFunny = false

	var c1 complex128 = 0.5i
	c := 10 + 3i

	m1, s1 := 10, "asd"

	x1 := 10
	f1 := float32(x1)

	var x = 10
	x2 := 20
	const a int = 10

	const myInt = 15

	numMessagesFromDoris := 72
	costPerMessage := .02
	totalCost := costPerMessage + float64(numMessagesFromDoris)

	mileage, company := 80276, "Toyota"

	helloLen := len("Hello World")
	helloLen1 := func() string { return "hellow" }

}
