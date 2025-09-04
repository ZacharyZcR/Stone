// Simple HTTP proxy benchmark
// Build with: go build -o benchmark benchmark_http.go
// Usage: ./benchmark -target=localhost:8080 -requests=1000

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	var (
		target   = flag.String("target", "localhost:8080", "WAF target address")
		requests = flag.Int("requests", 1000, "Number of requests")
		workers  = flag.Int("workers", 10, "Number of concurrent workers")
	)
	flag.Parse()

	fmt.Printf("Benchmarking Stone WAF HTTP handler\n")
	fmt.Printf("Target: %s\n", *target)
	fmt.Printf("Requests: %d\n", *requests)
	fmt.Printf("Workers: %d\n\n", *workers)

	start := time.Now()
	
	var wg sync.WaitGroup
	requestChan := make(chan int, *requests)
	
	// Fill request channel
	for i := 0; i < *requests; i++ {
		requestChan <- i
	}
	close(requestChan)
	
	// Start workers
	var successCount, errorCount int64
	var mu sync.Mutex
	
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{Timeout: 5 * time.Second}
			
			for range requestChan {
				resp, err := client.Get(fmt.Sprintf("http://%s/test", *target))
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					continue
				}
				
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}
	
	wg.Wait()
	elapsed := time.Since(start)
	
	fmt.Printf("Results:\n")
	fmt.Printf("Elapsed time: %v\n", elapsed)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Errors: %d\n", errorCount)
	fmt.Printf("Requests/sec: %.2f\n", float64(*requests)/elapsed.Seconds())
	fmt.Printf("Average latency: %.2fms\n", elapsed.Seconds()*1000/float64(*requests))
}