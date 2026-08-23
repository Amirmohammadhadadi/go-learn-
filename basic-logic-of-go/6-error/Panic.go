package __error

import "fmt"

func enrichUser(userID string) (User, error) {
	user, err := getUser(userID)
	if err != nil {
		// fmt.Errorf is GOATed: it wraps an error with additional context
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func enrichUser(userID string) User {
	user, err := getUser(userID)
	if err != nil {
		panic(err)
	}
	return user
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic:", r)
		}
	}()

	// this panics, but the defer/recover block catches it
	// a truly astonishingly bad way to handle errors
	enrichUser("123")
}
