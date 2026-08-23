package main

//func append(slice []Type, elems ...Type) []Type

func Append() {
	slice := make([]string, 10)
	slice = append(slice, "a")
	slice = append(slice, "b", "c", "d")
	slice = append(slice, []string{"a", "b", "c", "d"}...)
	//	dont do this!
	//someSlice := append(slice, "a")
}
