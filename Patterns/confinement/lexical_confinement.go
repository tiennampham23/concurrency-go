package confinement

import (
	"bytes"
	"fmt"
	"sync"
)

// Here we instantiate the channel within the lexical scope of the chanOwner function.
// This limits the scope of the write aspect of the results channel to the closure defined below it
// But channels are concurrent-safe
func LexicalConfinement() bool {
	wg := sync.WaitGroup{}
	wg.Add(1)
	chanOwner := func () <- chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <- chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}
	results := chanOwner()
	consumer(results)
	wg.Done()
	wg.Wait()
	return true
}

func LexicalConfinementWithBuffer(wg *sync.WaitGroup, data []byte) bool {
	defer wg.Done()
	var buff bytes.Buffer
	for _, b := range data {
		_, err := fmt.Fprintf(&buff, "%c", b)
		if err != nil {
			return false
		}
	}
	fmt.Println(buff.String())
	return false
}


func UseLexicalConfinementWithBuffer() {
	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go LexicalConfinementWithBuffer(&wg, data[:3])
	go LexicalConfinementWithBuffer(&wg, data[3:])
}