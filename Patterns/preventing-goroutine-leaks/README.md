We know goroutines are cheap and easy to crate, the runtime handles multiplexing the goroutines onto any
number of operating system thread so that we don't often have to worry about that level of
abstraction.

`But they do cost resources, and goroutines are not garbage collected by the runtime`
 
We must ensure a single child goroutine is guaranteed to be cleaned up.

The way to successfully mitigate this is to establish a signal between the parent goroutine and its children
that allows the parent to signal cancellation to its children.
By convention, this signal is usually a read-only channel named `done`. The parent goroutine passes this channel
to the child goroutine and then closes the channel when it wants to cancel the child goroutine

```go
// We pass the `done` channel to `doWork` function.  
func doWork(done <- chan interface{}, string <- chan string) {
	terminated := make(chan interface{})
	go func() {
	    defer fmt.Println("DoWork existed.")	
		defer close(terminated)
		for {
		    select {
			    case s := <- strings:
					// do something
                case <- done: 
					// Checking whether our done channel has been signaled. If it has, we return from the goroutine.
					return
            }           	
        }       
    }   
	return terminated
}

done := make(chan interface{})
terminated := DoWork(done, nil)
// Create another goroutine that will cancel the goroutine spawned in `doWork` if more than one second passes
go func() {
	// Cancel the operation after 1 second
	time.Sleep(1 * time.Second)
	fmt.Println("Canceling doWork goroutine")
	close(done)
}()
// This is where we join the goroutine spawned from `doWork` with the main goroutine
<- terminated
fmt.Println("Done")

// You can see that despite passing in nil for our `strings` channel, our goroutine still exits successfully
```


`If a goroutine is responsible for creating a goroutine, it is also responsible for ensuring it can stop the goroutine`