package __base

import "fmt"

func main() {
	fmt.Printf("Hello World\n")

	fmt.Errorf("error")

	fmt.Printf("%+v\n", struct{}{})
}
