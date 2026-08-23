package channel

import (
	"fmt"
	"time"
)

func CCMain() {
	ch := make(chan int)
	//	do some stuff with the channel
	close(ch)
	//v, ok := <-ch

}

// example
func sendReports(numBatches int, ch chan int) {
	for i := 0; i < numBatches; i++ {
		numReports := i*23 + 32%17
		ch <- numReports
		fmt.Printf("Send batch of %v reports \n", numReports, i)
		time.Sleep(time.Millisecond * 100)
	}
	close(ch)
}

func CCTest(numBatches int) {
	numSendCh := make(chan int)
	go sendReports(numBatches, numSendCh)
	fmt.Println("Start countin...")
	numReports := countReports(numSendCh)
	fmt.Printf("%v reports sent!\n", numReports)
}

func countReports(numSentCh chan int) int {
	total := 0
	for {
		numSent, ok := <-numSentCh
		if !ok {
			break
		}
		total += numSent
	}
	return total
}
