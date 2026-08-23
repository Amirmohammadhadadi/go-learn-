package main

import "fmt"

/*
arrayname[lowIndex:highIndex]
arrayname[lowIndex:]
arrayname[:highIndex]
arrayname[:]
*/

func array() {
	var _ [10]int
	var x [10]int = [10]int{
		1,2,
	}
	primes := [6]int{2, 3, 5, 7, 11, 13}
	mySlice := primes[1:4]
	// mySlice = {3, 5, 7}
	fmt.Printf("myInts: %v\n", mySlice)
}
