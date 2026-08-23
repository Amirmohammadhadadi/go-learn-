package channel

import "fmt"

func main() {
	var c chan string //c is nil
	c <- "hello"      //blocks

	fmt.Println(<-c) //blocks

	var d = make(chan int, 100)
	close(d)
	d <- 1           //panic: send on closed channel
	fmt.Println(<-d) //0
}
