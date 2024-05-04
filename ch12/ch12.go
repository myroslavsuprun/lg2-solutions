package main

import (
	"fmt"
	"math"
	"sync"
)

func main() {
	// first()
	// second()
	third()

}

var cacheMap = sync.OnceValue(newMap)

// 3.
// Write a function that builds a map[int]float64 where the keys are the numbers from 0 (inclusive) to 100,000 (exclusive)
// and the values are the square roots of those numbers (use the math.Sqrt function to calculate square roots).
// Use sync.OnceValue to generate a function that caches the map returned by this function
// and use the cached value to look up square roots for every 1,000th number from 0 to 100,000.
func third() {
	for i := 0; i < 100_000; i += 1000 {
		fmt.Printf("%v: %v \n", i, cacheMap()[i])
	}
}

func newMap() map[int]float64 {
	data := make(map[int]float64, 100_000)

	for i := 0; i < 100_000; i++ {
		data[i] = math.Sqrt(float64(i))
	}

	return data

}

// 2.
// Create a function that launches two goroutines.
// Each goroutine writes 10 numbers to its own channel.
// Use a for-select loop to read from both channels, printing out the number and the goroutine that wrote the value.
// Make sure that your function exits after all values are read and that none of your goroutines leak.
func second() {
	fmt.Println("Second")

	ch1 := make(chan int, 10)
	go func() {
		defer close(ch1)
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
	}()

	ch2 := make(chan int, 10)
	go func() {
		defer close(ch2)
		for i := 0; i < 10; i++ {
			ch2 <- i
		}
	}()

	for count := 0; count < 2; {
		select {
		case v, ok := <-ch1:
			if !ok {
				ch1 = nil
				count++
				continue
			}
			fmt.Println("ch1", v)
		case v, ok := <-ch2:
			if !ok {
				ch2 = nil
				count++
				continue
			}
			fmt.Println("ch2", v)
		}
	}
}

// 1.
// Create a function that launches three goroutines that communicate using a channel.
// The first two goroutines each write 10 numbers to the channel.
// The third goroutine reads all the numbers from the channel and prints them out.
// The function should exit when all values have been printed out.
// Make sure that none of the goroutines leak. You can create additional goroutines if needed.
func first() {
	fmt.Println("First")

	ch := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	go func() {
		for i := 0; i < 4_000_000; i++ {
			fmt.Println(i)
			ch <- i + 10
		}
	}()

	go func() {
		defer wg.Done()
		nms := make([]int, 00, 20)
		for len(nms) <= 20 {
			select {
			case v := <-ch:
				nms = append(nms, v)
			}
		}
		fmt.Println(nms)

	}()

	wg.Wait()

}
