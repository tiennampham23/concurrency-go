## The or-channel

> All times you may find yourself wanting to combine one or more `done` channels into a single `done` channel that closes if any of its
> components channels close.
> It is perfectly acceptable, to write a `select` statement that performs this coupling
> However, sometimes you can't know the number of `done` channel you're working with at runtime.


```go
var or func(channels ...<-chan interface{}) <- chan interface{}

or = func(channels ...<- chan interface{}) <- chan interface{} {
	switch len(channels) {
	    case 0:
			return nil
        case 1:
			return channels[0]
    }
	orDone := make(chan interface{})
	go func() {
	    defer close(orDone)	
		switch len(channels) {
		    case 2:
				select {
				    case <- channels[0]:
                    case <- channels[1]:
                }   
				default:
					select {
					case <- channels[0]:
                    case <- channels[1]:
                    case <- channels[2]:
                    case <- or(append(channels[3:], orDone)...)
                    }   
        }      
    }()
	return orDone
}

sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
    }()
	return c
}

start := time.Now()
<- or(
	    sig(2 * time.Hour)
        sig(5 * time.Minute)
        sig(1 * time.Second)
        sig(1 * time.Hour)
        sig(1 * time.Minute)
)
fmt.Printf("Done after: %v", time.Since(start))
```

This pattern is useful to employ at the intersection of modules in your system.
At these intersections, you tend to have multiple conditions for canceling trees of goroutines through your call stack.
Using the `or` function you can simply combine these together and pass it down the stack.