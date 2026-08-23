package __struct

type car struct {
	brand   string
	model   string
	doors   int
	mileage int
	// wheel is a field containing an anonymous struct
	wheel struct {
		radius   int
		material string
	}
}

var myCar = car{
	brand:   "Rezvani",
	model:   "Vengeance",
	doors:   4,
	mileage: 35000,
	wheel: struct {
		radius   int
		material string
	}{
		radius:   35,
		material: "alloy",
	},
}

func _main() {
	myCar := struct {
		brand string
		model string
	}{
		brand: "Toyota",
		model: "Camry",
	}
}
