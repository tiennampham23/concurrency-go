In concurrent programs, error handling can be difficult to get right. We spend so much time thinkking about
how our various processes will be sharing information and coordinating, we forget to consider
how they'll gracefully handling errored states.

**Who should be responsible for handling the error?**

Take a loot at the following example:
```go
checkStatus := func(done <- chan interface, urls ...string) <- chan *http.Response {
    responses := make(chan *http.Response)	
	go func() {
		defer close(responses)
		for _, url := range urls {
	        resp, err := http.Get(url)		
			if err != nil {
			    fmt.Println(err)	
				continue
            }   
			select {
			case <- done:
				return
            case responses <- resp:
            }
        }
    }
	return responses
}


done := make(chan interface{})
defer close(done)
urls := []string{"https://google.com", "https://facebbook.com"}
for response := range checkStatus(done, urls) {
	fmt.Printf("Response: %v", response.Status)
}
```

Here you see that the goroutine has been given no choice in the matter. It can't simply swallow the error,
so it does the only sensible thing: it only prints the error.

The following example demonstrates a correct solution:
```go
type Result struct {
	Error error
	Response *http.Response
}
checkStatus := func(done <- chan interface{}, urls ...string) <- chan Result {
	results := make(chan Result)
	go func() {
	    defer close(results)	
		for _, url := range urls {
		    var result Result
			resp, err := http.Get(url)
			result = Result {
			    Error: err,
				Response: resp,
            }   
			select {
			case <- done:
			    return
            case results <- result:
            }  
        }
    }()
}

done := make(chan interface{})
defer close(done)

urls := []string{"https://google.com", "https://facebook.com"}
for result := range checkStatus(done, urls...) {
	if result.Err != nil {
	    fmt.Printf("err %v", result.Error)	
    }
	fmt.Printf("Response %v", result.Response.Status)
}
```