package main

import "fmt"

/*
for INITIAL; CONDITION; AFTER{
// do something
}

	for INITIAL; ; AFTER {
	  // do something forever
	}

	for CONDITION {
	  // do some stuff while CONDITION is true
	}

	for {
	  // do some stuff forever
	}
*/
func _main() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	// Prints 0 through 9

	//while
	plantHeight := 1
	for plantHeight < 5 {
		if plantHeight == 4 {
			plantHeight = 0
			break
		}
		fmt.Println("still growing! current height:", plantHeight)
		plantHeight++
	}
	fmt.Println("plant has grown to ", plantHeight, "inches")
}

//continue
//break
