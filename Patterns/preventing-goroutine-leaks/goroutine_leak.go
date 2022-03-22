package main

import (
	"fmt"
	"math/rand"
	"time"
)

func doWork(strings <- chan string) <- chan interface{} {
	completed := make(chan interface{})
	go func() {
		defer fmt.Println("Do work existed")
		defer close(completed)
		for s := range strings {
			fmt.Println(s)
		}
	}()
	return completed
}

func newRandStream(done <-chan interface{}) <-chan int {
	randStream := make(chan int)
	go func() {
		defer fmt.Println("newRandStream closure exited.")
		defer close(randStream)
		for {
			select {
			case randStream <- rand.Int():
			case <-done:
				return
			}
		}
	}()

	return randStream
}

func main() {
	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	// Simulate ongoing work
	time.Sleep(1 * time.Millisecond)
}