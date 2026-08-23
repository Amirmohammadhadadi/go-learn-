package main

func base() {

	arr := make([]int, 0)
	arr = add1(arr)
	add2(&arr)

	//map1:=make(map[string]*int)

	//var x *int
	//fmt.Println(*x)
	//
	//var u int
	//
	//
	//var p *int
	//p = new(int)
	//*p = 10
	//p = nil
	//fmt.Println(*p + 10)
	//
	//
	////var p *int
	//// do not make pointer nil
	//
	//myString := "hello"
	//myStringPtr := &myString
	//
	//fmt.Println(*myStringPtr) //read myString through the pointer
	//*myStringPtr = "world"    //set myString through the pointer
}

func add1(arr []int) []int {
	arr = append(arr, 1)
	return arr
}
func add2(arr *[]int) {
	*arr = append(*arr, 2)
}
