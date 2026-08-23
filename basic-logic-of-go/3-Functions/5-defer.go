package __Functions

import (
	"fmt"
	"os"
)

func GetUsername(dstName, srcName string) (username string, err error) {
	// Open a connection to a database
	conn, _ := db.Open(srcName)

	// Close the connection *anywhere* the GetUsername function returns
	defer conn.Close()

	username, err = db.FetchUser()
	if err != nil {
		// The defer statement is auto-executed if we return here
		return "", err
	}

	// The defer statement is auto-executed if we return here
	return username, nil

}

/*
Multiple Defers

	The location of a defer statement inside a function matters.
	The deferred call is registered at the point where defer is executed,
	and it will run when the function returns.
	If you have multiple defer statements in a single function,
	they are executed in last-in,
	first-out order (the last deferred call runs first).
*/
func CreateTempFile() {
	f, _ := os.Create("temp-42.txt")
	defer os.Remove(f.Name()) // executed second
	defer f.Close()           // executed first

	fmt.Fprintln(f, "How many roads must a man walk down?")
}
