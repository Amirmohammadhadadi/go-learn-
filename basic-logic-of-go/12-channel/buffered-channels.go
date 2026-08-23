package channel

import "fmt"

func BCMain() {
	//ch := make(chan int,100)
}

//example

func addEmailsToQueue(emails []string) chan string {
	emailsToSend := make(chan string, len(emails))
	for _, _email := range emails {
		emailsToSend <- _email
	}
	return emailsToSend
}

func sendEmails(batchSize int, ch chan string) {
	for i := 0; i < batchSize; i++ {
		_email := <-ch
		fmt.Println("Sending email:", _email)
	}
}
func BCtest(emails ...string) {
	fmt.Printf("adding %v emails to queue...\n", len(emails))
	ch := addEmailsToQueue(emails)
	fmt.Println("Sending emails...")
	sendEmails(len(emails), ch)
}
