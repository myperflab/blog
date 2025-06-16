---
title: "Go Code Profiling Guide with Examples"
thumbnailImagePosition: left
thumbnailImage: /images/s1.jpeg
date: 2024-12-09
categories:
- perf engg
- perf testing
tags:
- Go Lang
showPagination: true
draft: False
---
<!--more-->

A perf engineer's best companions are profiling tools which can help them to reveal the bottlenecs of the application. There are a variety of the tools\libs are available. In this write-up we will discuss about how golang libs\tools can be used and how much effective + powerfull they are at capturing the profiles of the code\application which can be used for analysis to identify the bottlenecks. 

Go Provides built-in tools for profiling CPU, memory, goroutines, and more through the pprof and runtime packages. This guide walks through profiling a sample Go program with performance issues and demonstrates how to detect and address them.

# Types of Profiling in Go

| Profile Type | Description |
|--------------|-------------|
| CPU Profile | Shows where the CPU time is spent |
| Memory (Heap) | Tracks memory allocations and leaks |
| Goroutine Dump | Lists current goroutines and stack traces |
| Block Profile | Reveals blocking operations like channel ops |
| Mutex Profile | Identifies lock contention |
| Execution Trace | Provides detailed timeline of program events |


## Sample Go Program with CPU and Memory Issues

💡 Make sure to have **Go language** installed to run these experiments and collect the profiling dumps (**pprofs**).

```go
package main

import (
    "fmt"
    "math"
    "net/http"
    _ "net/http/pprof"  // Import pprof handlers
    "sync"
)

// Worker function that performs CPU-intensive calculations
func worker(id int, jobs <-chan int, results chan<- float64, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        // Simulate CPU-intensive work
        result := 0.0
        for i := 0; i < 1000000; i++ {
            result += math.Sqrt(float64(i + j))
        }
        results <- result
    }
}

func processJobs() {
    // Number of workers and jobs
    numWorkers := 4
    numJobs := 100

    // Create channels
    jobs := make(chan int, numJobs)
    results := make(chan float64, numJobs)

    // Create WaitGroup to wait for all workers to finish
    var wg sync.WaitGroup

    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    // Start a goroutine to collect results
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect and print results
    sum := 0.0
    for result := range results {
        sum += result
    }

    fmt.Printf("Total sum: %f\n", sum)
}

func main() {
    // Start HTTP server with pprof handlers
    go func() {
         if err := http.ListenAndServe("localhost:4567", nil); err != nil {
            fmt.Printf("Error starting server: %v\n", err)
        }
    }()

    // Run the workload continuously
    for {
        processJobs()
    }
}
```

This example demonstrates several key concepts:

1. **HTTP-based Profiling**: Uses Go's built-in pprof HTTP server
2. **Concurrent Processing**: Uses multiple goroutines to process jobs in parallel
3. **Channel Communication**: Uses channels for job distribution and result collection
4. **Synchronization**: Implements WaitGroup to ensure all goroutines complete

To run and analyze this program:

1. **Save** the code to a file (e.g., `cpu_profile_demo.go`)
2. **Run** the program:
   ```bash
   go run cpu_profile_demo.go
   ```
3. In a separate terminal, collect profiles using:
   ```bash
   # For CPU profile (30 seconds)
   go tool pprof http://localhost:4567/debug/pprof/profile?seconds=30
   
   # For heap profile
   go tool pprof http://localhost:4567/debug/pprof/heap
   
   # For goroutine profile
   go tool pprof http://localhost:4567/debug/pprof/goroutine

   or, simply run the following in any browser or cur
   http://localhost:4567/debug/pprof/goroutine?debug=1
   http://localhost:4567/debug/pprof/goroutine?debug=2
   http://localhost:4567/debug/pprof/goroutine?debug=3

     ```
    debug levels - 1 or 2 or 3 defines the depth of the collected information. In my POV depth level 2 is good enough to identify the blocked or problamatic go routines. 

    for the reference, attaching some sample pprof files

    📁 [Download sample pprof files](/images/pprof.zip) 

4. After collecting the pprofs, next things comes up - how to open/view?

    **Note:** to render the graphs - Graphviz is required. for linux: apt install graphviz, for windows dpwnload it from internet and install it.
    ```bash
    #pprofs:
    go tool pprof -http=:8080 <pprof file name>
    use any free port, in my case 8080 is free.
    it will open a page in the browser, in the view tab, there are multiple options to view the info. Use as fits best.

    #Dumps collected over brouser
    Use any text editor.
    ```
    for the code example, we can see following info from pprofs

    a. **CPU pprofs**

    ![](/images/cpu_pprof.png) 
 
    b. Similarly other pprofos (**mem and goroutines**) can be examined.

    c. **goroutine dumps** in text files:
    if we check, goroutine_debug_2.txt (in the zip file)
    - there are 4 goroutines as per code, same can be seen in the dumos and there respective stack trace.
    - go routine are similar to threads in java
    - there are more information about channel and other go routines

5. What information can be inferred from these pprofs or dumps:
    The profile will show you:
    - Which goroutines are consuming the most CPU time
    - How the work is distributed across goroutines
    - Any potential bottlenecks in the concurrent processing

    This example is particularly useful for:
    - On-demand profiling without modifying the code
    - Collecting different types of profiles (CPU, heap, goroutines)
    - Profiling long-running applications
    - Real-time performance analysis
SO far, we have covered CPU, Mem and go routine profiling. Now lets jump into Block Profiling 

## Block Profiling

As the name suggest, A Block Profile captures where goroutines are blocked, i.e., waiting for synchronization events like:
- Waiting to receive/send on a channel
- Waiting to acquire a mutex (lock)
- Waiting on select, Cond, etc.

These are often the performance bottlenecks in concurrent programs. A block profile helps in identify -

- Where goroutines block in your code
- How long they are blocked
- How often blocking happens at a specific code location

How to capture block profile: similar as CPU profiling, only diff is end point

```bash
go tool pprof http://localhost:4567/debug/pprof/block?seconds=30
```
### use following code block to simulate
```bash
package main

import (
	"fmt"
	"math"
	"net/http"
	_ "net/http/pprof" // Import pprof handlers
	"sync"
	"time"
)
// demonstrateBlockProfiling creates scenarios that will show up in block profiling
func demonstrateBlockProfiling() {
	// Create channels with very small buffers to force more blocking
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	// Create multiple mutexes for more contention
	var mu1, mu2 sync.Mutex

	// Launch more goroutines with longer running times
	for i := 0; i < 10; i++ {
		go func(id int) {
			fmt.Printf("Starting blocking goroutine %d\n", id)
			for {
				// Contend for first mutex
				mu1.Lock()
				// Simulate some work
				time.Sleep(time.Millisecond * 100) // Increased sleep time
				mu1.Unlock()

				// Contend for second mutex
				mu2.Lock()
				time.Sleep(time.Millisecond * 80) // Increased sleep time
				mu2.Unlock()

				// Force blocking on channel operations
				ch1 <- id // This will block when channel is full
				time.Sleep(time.Millisecond * 50)
				<-ch1

				ch2 <- id // This will block when channel is full
				time.Sleep(time.Millisecond * 50)
				<-ch2
			}
		}(i)
	}
}

func main() {
	// Enable block profiling
	runtime.SetBlockProfileRate(1) // Enable block profiling with rate 1 (every blocking event)

	// Start HTTP server with pprof handlers
	go func() {
		fmt.Println("Starting pprof server on http://localhost:4567")
		if err := http.ListenAndServe("localhost:4567", nil); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Start block profiling demonstration in background
	fmt.Println("Starting block profiling demonstration...")
	go demonstrateBlockProfiling()

}

```
**Note:** runtime.SetBlockProfileRate(1) it is required to have block profile.
- 0 (default) = block profiling is disabled
- 1 = record every blocking event
- n > 1 = record 1 in n blocking events
- negative = disable block profiling

Without calling runtime.SetBlockProfileRate(1), the block profiler won't collect any data. This is different from CPU profiling which is enabled by default when you import net/http/pprof. Block profiling is more expensive in terms of performance overhead, which is why it's opt-in rather than enabled by default.

pprof let's know which goroutine, loc or channel is locked while sampling. And that's how it is diff from normal CPU profiling, as cpu profilied tells which are the func/methods are consuming most of the cpu.

In general, CPU profiling is more common than block profiling. With this, quetion comes when should be block profiles is used:

Use it when:
- App is running slow but CPU isn't maxed
- suspect lock contention
- timeouts, hung handlers, or queue backups

It helps you find the bottlenecks not in computation, but in coordination between parts of your program.

## Mutex Profiling
Roughly, we can say mutex profiling is a sub set of Block Profiling, where it tracks down only **how long goroutines wait to acquire a mutex (specifically sync.Mutex or sync.RWMutex).**

similar to block profile, it has to be enabled with ***runtime.SetMutexProfileFraction(n)***

### How it works:
When multiple goroutines try to acquire the same mutex, the ones that have to wait are measured.
The profile records:
- Lock wait time
- Number of blocking instances
- Stack trace where it happened

