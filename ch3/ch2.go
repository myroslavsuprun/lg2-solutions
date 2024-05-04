package ch2

import (
	"fmt"
)

func main() {
	// 1
	greetings := []string{"Hello", "Hola", "नमस्कार", "こんにちは", "Привіт"}

	f := greetings[:2]
	s := greetings[1:4]
	t := greetings[3:5]

	fmt.Println(greetings)
	fmt.Println(f, s, t)

	// 2
	message := []rune("Hi 👩 and 👨")

	r := string(message[3])
	fmt.Println(r)

	// 3
	type Employee struct {
		firstName string
		lastName  string
		id        int
	}

	first := Employee{
		"Myroslav",
		"Suprun",
		0,
	}

	second := Employee{
		firstName: "Myroslav",
		lastName:  "Suprun",
		id:        1,
	}

	var third Employee
	third.firstName = "Myroslav"
	third.lastName = "Suprun"
	third.id = 1

	fmt.Println(first, second, third)

}
