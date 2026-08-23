package channel

import (
	"fmt"
	"time"
)

// go doSomething()

// example
func sendEmail(message string) {
	go func() {
		time.Sleep(250 * time.Millisecond)
		fmt.Printf("Email received: %s\n", message)
	}()
	fmt.Printf("Email sent: %s\n", message)
}
