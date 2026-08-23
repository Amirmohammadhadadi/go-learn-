package __struct

import "fmt"

type Car struct {
	brand      string
	model      string
	doors      int
	mileage    int
	frontWheel Wheel
	backWheel  Wheel
}
type Wheel struct {
	radius   int
	material string
}

func _main() {
	myCar := Car{
		brand: "asd",
		frontWheel: Wheel{
			radius:   10,
			material: "strawberry",
		},
		doors:   2,
		mileage: 3,
		model:   "strawberry",
	}
	fmt.Println(myCar)
	fmt.Println(myCar.brand)
	fmt.Println(myCar.frontWheel.radius)

}
