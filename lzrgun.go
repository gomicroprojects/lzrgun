package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	// Importing the flag package (http://golang.org/pkg/flag/)
	"flag"
)

// we are declaring the vars here, which will be populated through the flag packages
var (
	// The number of total requests to perform
	totalRequests int

	// The number of concurrent requests to perform
	concurrentRequests int
)

// A sync.WaintGroup so we can wait for all workers to finish before exiting
// see http://golang.org/pkg/sync/#WaitGroup for an example of the WaitGroup
var working sync.WaitGroup

// The channel through which the workers will receive their workload
var workload chan Workload

var httpClient *http.Client

// Workload is a structure which descibes the action a worker should perform
// Currently it only contains an URL to perform an HTTP GET against, but this can easily be expanded
type Workload struct {
	URL *url.URL
}

// The main func. The execution of the program starts here
func main() {
	// See http://golang.org/pkg/flag/#IntVar
	// It takes a pointer to an int, we use the "&" operator to take the address of our int
	flag.IntVar(&totalRequests, "n", 100, "Number of total requests to perform.")
	flag.IntVar(&concurrentRequests, "nc", 10, "Number of concurrent requests to perform.")
	flag.Parse()

	// The number of arguments after all flags should be 1
	// We want the URL to be the sole (required) argument
	if flag.NArg() == 0 {
		fmt.Println("Requiring exactly one argument.")
		// see the function declaration below
		printUsage()
		os.Exit(1)
	}

	targetURL, err := url.Parse(flag.Arg(0))
	if err != nil {
		fmt.Printf("Error parsing URL: %s\n", err)
		printUsage()
		os.Exit(1)
	}
	fmt.Printf("using %s with %d total requests (%d concurrently)\n", targetURL.String(), totalRequests, concurrentRequests)

	// we initialize the workload channel and set the buffer size to 1.5 times the number of workers
	// so the work queue is nicely filled
	workload = make(chan Workload, concurrentRequests+(concurrentRequests/2))

	httpClient = &http.Client{}

	for i := 0; i < concurrentRequests; i++ {
		working.Add(1)
		go worker()
	}

	for i := 0; i < totalRequests; i++ {
		workload <- Workload{URL: targetURL}
	}
	close(workload)
	working.Wait()

}

func printUsage() {
	// os.Args[0] is the command used to start the binary
	fmt.Printf("Usage: %s [-n total] [-nc concurrent] url\n", os.Args[0])
	flag.PrintDefaults()
}

func worker() {
	for {
		workload, ok := <-workload
		if !ok {
			working.Done()
			return
		}
		_ = workload
	}
}
