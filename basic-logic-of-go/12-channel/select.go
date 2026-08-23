package channel

import "fmt"

func SMain() {
	//select {
	//case i, ok := <-chInts:
	//	fmt.Println(i)
	//case s, ok := <-chStrings:
	//	fmt.Println(s)
	//default:
	//	//	receiving from ch would block
	//	// so do something else
	//}

}

//example

func logEmail(email string) {
	fmt.Println("Email:", email)
}
func logSms(sms string) {
	fmt.Println("SMS:", sms)
}
func logMessages(chEmails, chSms chan string) {
	for {
		select {
		case _email, ok := <-chEmails:
			if !ok {
				return
			}
			logEmail(_email)
			break
		case _sms, ok := <-chSms:
			if !ok {
				return
			}
			logSms(_sms)
			break
		}
	}
}
