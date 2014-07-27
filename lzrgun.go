package main

import (
	"fmt"
	"net/url"
	"os"
	// Importing the flag package (http://golang.org/pkg/flag/)
	"flag"
)

// we are declaring the vars here, which will be populated through the flag packages
// The number of total requests to perform
var totalRequests int

// The number of concurrent requests to perform
var concurrentRequests int

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
}

func printUsage() {
	// os.Args[0] is the command used to start the binary
	fmt.Printf("Usage: %s [-n total] [-nc concurrent] url\n", os.Args[0])
	flag.PrintDefaults()
}
