package __error

import (
	"errors"
	"fmt"
	"strconv"
)

//type error interface {
//	Error() string
//}

func Atoi(s string) (int, error) {
	return 10, errors.New("wrong")
}

func _main() {
	// Atoi converts a stringified number to an integer
	i, err := strconv.Atoi("42b")
	if err != nil {
		fmt.Println("couldn't convert:", err)
		// because "42b" isn't a valid integer, we print:
		// couldn't convert: strconv.Atoi: parsing "42b": invalid syntax
		// Note:
		// 'parsing "42b": invalid syntax' is returned by the .Error() method
		return
	}
	// if we get here, then the
	// 1-variable i was converted successfully
	fmt.Println(i)
}

type userError struct {
	name string
}

func (e userError) Error() string {
	return fmt.Sprintf("%v has a problem with their account", e.name)
}
func canSendToUser(userName string) bool {
	return false
}
func sendSMS(msg, userName string) error {
	if !canSendToUser(userName) {
		return userError{name: userName}
	}
	var err error = errors.New("something went wrong")
	return err
}
