package main

import (
	"fmt"
	"time"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func MakePerson(firstName string, lastName string, age int) Person {
	return Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

func MakePersonPointer(firstName string, lastName string, age int) *Person {
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}

}

func UpdateSlice(arr []string, val string) {
	arr[len(arr)-1] = val
	fmt.Println(arr)
}

func GrowSlice(arr []string, val string) {
	arr = append(arr, val)
	fmt.Println(arr)
}

func main() {
	sl := []string{"1", "2", "3"}

	fmt.Println(sl)
	UpdateSlice(sl, "4")

	GrowSlice(sl, "4")

	manyPersons()
}

const PEOPLE = 10_000_000

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func manyPersons() {
	defer timer("manyPersons")()
	// var persons []Person = make([]Person, PEOPLE, PEOPLE)
	var persons []Person
	for i := 0; i < PEOPLE; i++ {
		// persons[i] = MakePerson("Myro", "Jagaja", i+1)
		persons = append(persons, MakePerson("Myro", "Jagaja", i+1))
	}

}
