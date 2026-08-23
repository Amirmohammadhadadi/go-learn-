package __variable

import "fmt"

func mainString() {
	fmt.Println("Hello world")
	ss := fmt.Sprintf("abc %v", 1)
	fmt.Println(ss)
}
