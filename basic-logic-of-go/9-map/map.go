package main

import (
	"fmt"
)

func _map() {
	ages := make(map[string]int)
	ages["gophers1"] = 42
	ages["gophers2"] = 43
	ages["gophers1"] = 44 //overwrites 42
	ages1 := map[string]int{
		"gophers1": 5,
	}
	fmt.Println(len(ages1))
	if _, ok := ages["gophers1"]; ok {
		delete(ages, "gophers1")
	}
	//if ages1['gophers1'] {
	//	delete(ages1, "gophers1")
	//}
}

