package __struct

import "fmt"

type car struct {
	brand string
	model string
}

type truck struct {
	// "car" is embedded, so the definition of a
	// "truck" now also additionally contains all
	// of the fields of the car struct
	car
	bedSize int
}

func _main() {
	lanesTruck := truck{
		bedSize: 10,
		car: car{
			brand: "Toyota",
			model: "Tundra",
		},
	}

	fmt.Println(lanesTruck.brand) // Toyota
	fmt.Println(lanesTruck.model) // Tundra
}
