package channel

import (
	"fmt"
	"time"
)

func RMain() {
	ch := make(chan string)
	for item := range ch {
		println(item)
	}
}

// example
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
		time.Sleep(time.Microsecond * 10)
	}
	close(c)
}
func RTest(n int) {
	fmt.Printf("Printing %v numbers...\n", n)
	concurrentFib(n)
}
func concurrentFib(n int) {
	chInts := make(chan int)
	go func() {
		fibonacci(n, chInts)

	}()
	for item := range chInts {
		fmt.Println(item)
	}
}
