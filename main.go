package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

func main() {

	mode := flag.String("mode", "", "Mode to run: queue | producer | consumer")
	flag.Parse()

	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	switch *mode {
	case "queue":
		declareQueue()
	case "producer":
		runProducer(100) // Default to 100 messages for producer
	case "consumer":
		runConsumer(10, 100, 100) // Default to 10 workers, 100 messages, and 10ms delay
	case "benchmark":
		RunBenchmark()
	default:
		fmt.Println("Please specify --mode=[queue|producer|consumer]")
		os.Exit(1)
	}

	fmt.Printf("GOMAXPROCS set to %d\n", n)

}
