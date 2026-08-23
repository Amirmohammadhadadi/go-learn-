package main

import "fmt"

func functions() {
	// func make([]T, len, cap) []T
	mySlice := make([]int, 5, 10)

	// the capacity argument is usually omitted and defaults to the length
	mySlice := make([]int, 5)

	mySlice := []string{"I", "love", "go"}

	mySlice := []string{"I", "love", "go"}
	fmt.Println(len(mySlice)) // 3

	mySlice := []string{"I", "love", "go"}
	fmt.Println(cap(mySlice)) // 3

	mySlice := []string{"I", "love", "go"}
	fmt.Println(mySlice[2]) // go

	mySlice[0] = "you"
	fmt.Println(mySlice) // [you love go]
	var arr [10]string = [10]string{"1", "2", "3", "4", "5"}
	fmt.Println(len(arr))
	fmt.Println(cap(arr))
}

func getMessages(messages []string) []string {
	return messages
}

// deleteArrayByIndex(&arr, j)
func deleteArrayByIndex(arr *[]uint8, j int) {

	*arr = append((*arr)[:j], (*arr)[j+1:]...)
}
