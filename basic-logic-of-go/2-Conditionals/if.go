package __Conditionals

import "fmt"

func ifMain() {
	//if INITIAL_STATEMENT; CONDITION {
	//}

	//if statements in Go do not use parentheses around the condition:
	if height > 4 {
		fmt.Println("You are tall enough!")
	}

	//	else if and else are supported as you might expect:
	if height > 6 {
		fmt.Println("You are super tall!")
	} else if height > 4 {
		fmt.Println("You are tall enough!")
	} else {
		fmt.Println("You are not tall enough!")
	}

	length := getLength(email)
	if length < 10 {
		fmt.Printf("Email must be at least 10 characters, is %d\n", length)
	}
	if length := getLength(email); length < 10 {
		fmt.Printf("Email must be at least 10 characters, is %d\n", length)
	}

	if val, ok := l.Aparteman.(map[string]interface{}); ok {
		fmt.Println("این یک آبجکت است", val)
	}

}
