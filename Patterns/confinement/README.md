There are two ways to work with concurrent code:
* The first, Synchronization primitives for sharing memory, like using `sync.Mutex` or `Sync.RWMutex`
* The second, Synchronization via communicating, like `channel`

However, there are a lot of other options that are implicity safe within multiple concurrent processes:
* Immutable data: using the copy of data, not pointers
* Data protected by confinement: ensuring data is only ever available from one concurrent process.

There are two kinds of confinement possible: ad hoc and lexical
* `Ad hoc confinement`: is when you achieve confinement through a convention - whether it be set but the languages community,
the group you work within, or the codebase you work within. That's difficult in the reality, even those you have some tools to
analysis code before committing.
```go
data := make([]int, 4)
loopData := func(handleData chan <- int) {
	defer close(handleData)
	for i := range data {
		handleData <- data
    }
}
handleData := make(chan int)
go loopData(handleData)
for num := range handleData {
	fmt.Println(num)
}
// We can see that the `data` slice of integers is available from both the loopData function
// and the loop over the handleData channel;
// However, by convention we're only accessing it from the `loopData` function
// This can be wrong when some developers are faced deadline =))
```

* Lexical confinement: involves using lexical scope to expose only the correct data and concurrency primitives for multiple
concurrent processes to use.
It makes it impossible to do the wrong thing.
```go
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
```


