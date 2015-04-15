Wow, that's some interesting talk from Kai!

Hi everyone, my name is Mark and I am working at Viki.

Today, I am going to talk about the Go garbage collector. This is also a lightning talk, so it's gonna be short!
 
One of the attrative points of Go is its garbage collector. Everyone knows that it's garbage collecting unused memory, but do you know exactly why you need it and how it works, as well as how it might affect you?

This talk is an introduction to what you need to know about the GC.

First, I will go briefly to memory allocation in C, because it is the reason why garbage collector was invented in the first place. Then we will see how Go do it memory allocation. This will help you understand why you need the garbage collector.

After that, we will just see how it works, and how does it matter to you as a programmer
So some history. First, how many of you did C or C++ before ?

->

Great!

So C is a language where you have to manually allocate memory space for your variable, unlike many modern programming language like Ruby, Python, Go and even Java. When you allocate some memory in C, it either goes to a stack or a heap. These are just names for 2 different memory space on your computer with different behavior.

For example, to allocate memory on stack, you would do something similar to the program on the left. It's a very simple program to find max of 2 integer. In the max function, I declared a variable result and assign it some value. Then I return the value back to the caller.

This means the variable result is accessed within the function only. Hence it is allocated from the stack.

When the max function return, the stack space is destroyed, hence nobody else can use the value of result.

On the other hand, I can decide to write my max function in this way:

This time it returns a pointer to some integer, which I de-reference and print. To do that, I have to reserve a space somewhere that is not the stack, otherwise I won't be able to access the result.

Thus I used malloc here, and allocated some space in the heap. The space to hold the address itself is on the stack, however.

When I finished max, the space holding the address is destroyed, but the space holding the value still exists. And because I have the address, I can access the content to print.

Now, you will notice that I have to do a free at the end, otherwise nobody else can use the space I used on heap.

This manual allocation / deallocation can be very error prone, especially if you work with threads and concurrent things.

Being lazy, a programmer might have asked: is there a way for me not to care about freeing memory? An d thus garbage collection is born

Now, back to Go. I have a similar program.
Can someone tell me in this case, where is result allocated?


Right, in theory it's correct, execept that there is no official "stack" in Go, but the variable is allocated in the same way as in C: to a temporary space when the function run, and deallocated after the function exit

Now look at this example: We have return an average of 2 Gopher struct and return. If whatever we are allocating here is in stack, we wouldn't be able to access it after the breed fucntion return. So this will have to be allocated somewhere similar to the C heap. The difference here is that you don't need to do any deallocation.

That's because Go introduced a garbage collector. It's somewhat similar to the GC in Java, Ruby, Python and all, execpte that it's ran by tiny little Gophers running around buring unused memory.

How does that happen: Let's look at a demo

Here I have a HTTP server, where each request it will give you back one Gopher
Now, I will first compile and run it without the GC
Now, I will send a bunch of request to the server.
As you can see, it's starting to consume a lot memory. Let's look at its memory profile.
800 MB is not very good for a simple http server isn't it. This is because we turned the GC off.

Now if I enable it.
I will send the same benchmark test to the server this time
As you can see, the memory usage did increase, but not so crazy like last time

This time it only use less that 1MB of data, so the GC played a big difference here.

What actually happened is this
<Run the compiler with gcflags=-m>
You can see that the Gopher variable escape to heap, and it stays there
So every request you have there is one leaked space, and after a lot of request, this become huge (800MB)

After your request had finished, the GC detects that nobody need the newly created Gopher anymore and collect it from there

So, why should you care about the GC?
In an ideal world, the GC will just silently run and you will never notice it. However, Go 1.4 GC is not like that. Whenever it runs, it wouls stop every goroutine to do GC.

In Viki, we used to look at our error graph and saw timeout happen in this pattern.
After much looking, we found out that was the GC running, and when it we can't respond to client. That GC took around 2-3 seconds and many of our clients timed out then.

Another more scary scenario can be described this way:
You are writing a money transfer system from A to B, both in Go
...
now somebody got twice the money, and you are in trouble.

Last part: You can prevent these issue by understanding your memory usage and tuning your GC. There's already manything on the web so I won't repeat them here.

That's all for my presentation. Thank you