package __Interface

import (
	"fmt"
	"math"
)

type shape interface {
	area() float64
	perimeter() float64
	GetWidth() float64
}

type rect struct {
	width, height float64
}

func (r rect) GetWidth() float64 {
	return r.width
}
func (r rect) area() float64 {
	return r.width * r.height
}
func (r rect) perimeter() float64 {
	return 2*r.width + 2*r.height
}

type circle struct {
	radius float64
}

func (c circle) GetWidth() int {
	return 10
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c circle) perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func printShapeData(s shape) {
	fmt.Printf("Area: %v - Perimeter: %v\n", s.area(), s.perimeter())
}

func Test() {
	printShapeData(rect{3, 4})
	printShapeData(circle{1})
}
