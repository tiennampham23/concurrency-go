## Difference Between Concurrency and Parallelism

> Something that rún at the same time as something else. 
> Sometimes using the word `parallel` in this context is correct, but usually if the developers
> are discussing code, they really ought to be using the word `concurrent`.

Another:
> Concurrency is a property of the code, parallelism is a property of the running program.

Example: **If I write my code with the intent taht two chunks of the program will run in parallel, do I have any guarantee that will actually happen when the program is run?**.
**What happens if I run the code on a machine with only one core?**

=> The reveals a few interesting and important things:
* The first is that we don't write parallel code, only concurrent code 
that we hope will be run in parallel. Once again, parallelism is a property of the `runtime of our program, not the code`
* The second interesting thing is that we see it is possible - maybe even desirable - to be ignorant of whether our concurrent code is actually running in parallel.
The abstractions are what allow us to make the distinction between concurrency and parallelism. `Abtractions are the conccurency primitives, the program's runtime, the operating system, the platform the operating system runs on(container, VM.`
* The final interesting thing is that parallelism is a function of time, or context.
