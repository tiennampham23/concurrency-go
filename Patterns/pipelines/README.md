The following example is considered a pipeline stage

```go
multiply := func(values []int, multiplier int) []int {
	multipliedValues := make([]int, len(values))
	for i, v := range values {
	    multipliedValues[i] = v * multiplier	
    }
	return multiplierValues
}

// another stage

add := func (values []int, additive int) []int {
	addedValues := make([]int, len(values))
	for i, v := range values {
	    addValues[i] = v + additive
    }
	return addedValues
}
```


## Best practices for constructing pipelines
Channels are uniquely suited to constructing pipelines in Go.

```go
generator := func(done <- chan interace{}, integers ...int) <- chan int {
	intStream := make(chan int)
	go func() {
	    defer close(intStream)	
		for _, i := range integers {
		    select {
			case <- done:
				return
            case intStream <- i:
            }	
        }  
    }()
	return intStream
}

multiply := func(done <- chan interface{}, intStream <- chan int, multipler int) <- chan int {
	multipledStream := make(chan int)
	go func() {
	    defer close(multipledStream)	
		for i := range intStream {
		    select {
			case <- done:
				return
            case multipledStream <- i * multipler:
            }	
        }
    }()
	return multipliedStream
}

add := func(done <- chan interface{}, intStream <- chan int, additive int) <- chan int {
	addedStream := make(chan int)
    go func() {
        defer close(addedStream)
        for i := range intStream {
            select {
            case <- done:
                return
            case addedStream <- i + additive:
            }	
        }
    }()
    return addedStream
}

done := make(chan interface{})
defer close(done)
intStream :- generator(done, 1,2,3,4)
pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
for v := range pipeline {
	fmt.Println(v)
}
```

