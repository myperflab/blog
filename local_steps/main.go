package main

import (
	"fmt"
	"math"
	"net/http"
	_ "net/http/pprof" // Import pprof handlers
	"runtime"
	"sync"
	"time"
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
		fmt.Printf("Worker %d processed job %d\n", id, j)
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

	// Run the CPU workload
	for {
		processJobs()
	}
}
