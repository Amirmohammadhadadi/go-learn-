package main

import "fmt"

//for Index,Element := range SLICE{
//}

func _for() {
	fruits := []string{"apple", "banana", "orange", "grape"}

	for i, item := range fruits {
		fmt.Printf("fruits[%d] = %s\n", i, item)
	}
	length := len(fruits)
	if length > 1 {
		fmt.Printf("fruits[%d] = %s\n", length, fruits[length-1])
	}

}
