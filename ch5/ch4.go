package main

import (
	"fmt"
	"os"
)

func main() {
	// 1
	first()

	// 2
	flen, err := fileLen("ch4.go")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("File size is:", flen)

	// 3
	helloPrefix := prefixer("Hello")
	fmt.Println(helloPrefix("Bob"))
	fmt.Println(helloPrefix("Maria"))

}

func prefixer(prefix string) func(string) string {
	return func(v string) string {
		return fmt.Sprintf("%v %v", prefix, v)
	}
}

func fileLen(name string) (int, error) {
	f, closer, err := getFile(name)

	if err != nil {
		return 0, err
	}

	defer func() {
		closer()
	}()

	s, err := f.Stat()
	if err != nil {
		return 0, err
	}

	return int(s.Size()), nil
}

func getFile(name string) (*os.File, func(), error) {
	f, err := os.Open(name)

	if err != nil {
		return nil, nil, err
	}

	return f, func() {
		f.Close()
	}, err
}
