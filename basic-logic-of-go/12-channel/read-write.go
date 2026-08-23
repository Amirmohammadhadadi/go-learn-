package channel

import (
	"fmt"
	"time"
)

func RWMain() {
	ch := make(chan int)
	readCh(ch)
	writeCh(ch)
}

func readCh(ch <-chan int) {
	//	ch can only be read from
	//	in this function
	val := <-ch
	fmt.Println(val)
}
func writeCh(ch chan<- int) {
	//	ch can only be written from
	//	in this function
	ch <- 10
}

func saveBackups(snapshotTicker, saveAfter <-chan time.Time) {

	for {
		select {
		case <-snapshotTicker:
			fmt.Println("takeSnapshot")
		case <-saveAfter:
			fmt.Println("saveAfter")
		default:
			fmt.Println("waitForData")
			time.Sleep(time.Millisecond * 500)
		}
	}
}
