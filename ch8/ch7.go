package main

import "fmt"

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 |
		~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

func Double[T Integer](v T) T {
	return v * 2
}

func main() {
	// 1
	num := 22

	fmt.Println(num)
	num = Double(num)
	fmt.Println(num)

	// 2
	var pr Printing = 22

	DoWithPrintable(pr)

	// 3
	LLExec()
}

type Printable interface {
	~int | ~float64
	String() string
}

type Printing int

func (p Printing) String() string {
	return string(p)
}

func DoWithPrintable[T Printable](p T) {
	fmt.Println(p)
}
