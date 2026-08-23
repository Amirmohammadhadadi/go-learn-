package channel

import (
	"fmt"
	"time"
)

func mainChannels() {
	ch := make(chan int)
	ch <- 1
	v := <-ch
	fmt.Println("v", v)
	close(ch)
}

// example
type email struct {
	date time.Time
}

func filterOldEmails(emails []email) {
	isOldChan := make(chan bool)
	go func() {
		for _, e := range emails {
			if e.date.Before(time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC)) {
				isOldChan <- true
				continue
			}
			isOldChan <- false
		}
	}()
	isOld := <-isOldChan
	fmt.Println("email 1 is old:", isOld)
	isOld = <-isOldChan
	fmt.Println("email 2 is old:", isOld)
	isOld = <-isOldChan
	fmt.Println("email 3 is old:", isOld)
	// true,false,false
}

func waitForDbs(numDBs int, dbChan chan struct{}) {
	//answer
	for i := 0; i < numDBs; i++ {
		<-dbChan
	}
}
func getDatabaseChannel(numDBs int) chan struct{} {
	ch := make(chan struct{})
	go func() {
		for i := 0; i < numDBs; i++ {
			ch <- struct{}{}
			fmt.Printf("Database %v is online\n", i+1)
		}
	}()
	return ch
}
func Ctest(numDBs int) {
	dbChan := getDatabaseChannel(numDBs)
	waitForDbs(numDBs, dbChan)
}
