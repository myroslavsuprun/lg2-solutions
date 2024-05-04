package ch3

import (
	"fmt"
	"math/rand"
)

func main() {
	// 1
	rInts := make([]int, 100, 100)
	for i := 0; i < 100; i++ {
		rInts[i] = rand.Intn(101)
	}

	fmt.Println(rInts)

	// 2

	for _, v := range rInts {
		switch {
		case v%3 == 0 && v%2 == 0:
			fmt.Println("Six!")
		case v%2 == 0:
			fmt.Println("Two!")
		case v%3 == 0:
			fmt.Println("Three!")
		default:
			fmt.Println("Never mind")
		}
	}

	// 3
	var total int

	for i := 0; i < 10; i++ {
		total := total + i
		fmt.Println(total)
	}
	fmt.Printf("Total after the loop: %v\n", total)

}
