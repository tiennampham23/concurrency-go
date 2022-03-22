## Fan-Out, Fan-In

Fan out is a term to describe the process of starting multiple goroutines to handle input from
the pipeline. And fan-in is a term to describe the process of combining multiple results into
one channel.

You might consider fanning out one of your stages if both the following apply:
- It doesn't rely on values that the stage had calculated before
- It takes a long time to run

The property of order-independence is import because you have no guarantee in what
order concurrent copies of your stage will run, nor in what order they will return.

