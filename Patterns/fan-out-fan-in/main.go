package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

func toInt(done <- chan interface{}, valueStream <- chan int) <- chan int  {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for v := range valueStream {
			select {
			case <- done:
				return
			case intStream <- v:
			}
		}
	}()
	return intStream
}

func repeatFn(done <- chan interface{}, fn func() int) <- chan int {
	valueStream := make(chan int)
	go func() {
		defer close(valueStream)
		for {
			select {
			case <- done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}


func take(done <- chan interface{}, valueStream <- chan int, num int) <- chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <- done:
				return
			case v := <-valueStream:
				takeStream <- v
			}
		}
	}()
	return takeStream
}

func primeFinder(done <- chan interface{}, randIntStream <- chan int) <- chan int {
	primeStream := make(chan int)
	go func() {
		defer close(primeStream)
		for {
			select {
			case <- done:
				return
			case v := <- randIntStream:
				if big.NewInt(int64(v)).ProbablyPrime(0) {
					primeStream <- v
				}
			}
		}
	}()
	return primeStream
}

func fanIn(done <- chan interface{}, channels ...<- chan int) <- chan int {
	var wg sync.WaitGroup
	multiplexedStream := make(chan int)
	multiplex := func(c <- chan int) {
		defer wg.Done()
		for i := range c {
			select {
			case <- done:
				return
			case multiplexedStream <- i:
			}
		}
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()
	return multiplexedStream
}

func main()  {
	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	r := func() int {return rand.Intn(50000000)}
	randIntStream := toInt(done, repeatFn(done, r))
	numFinders := 4
	fmt.Printf("Spining up %d prime finders.\n", numFinders)
	finders := make([]<- chan int, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t %v\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}
