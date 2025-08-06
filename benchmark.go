package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"
)

type Config struct {
	NumMessages int
	NumWorkers  int
	DelayMs     int // Simulated delay per worker
}

type BenchmarkResult struct {
	NumMessages int
	NumWorkers  int
	DelayMs     int
	Duration    time.Duration
	NumGC       uint32 // number of garbage collections
}

func createProfileFile(profileDir, profileType string, numWorkers int) (*os.File, error) {
	profilePath := filepath.Join(profileDir, fmt.Sprintf("%s-%dworkers.pprof", profileType, numWorkers))
	file, err := os.Create(profilePath)
	if err != nil {
		return nil, fmt.Errorf("could not create %s profile: %w", profileType, err)
	}
	return file, nil
}

func RunBenchmark() {
	declareQueue()

	profileDir := "profiles"
	if err := os.MkdirAll(profileDir, os.ModePerm); err != nil {
		log.Fatalf("could not create profile directory: %v", err)
	}

	tests := []Config{
		{NumMessages: 1000, NumWorkers: 5, DelayMs: 100},
		{NumMessages: 1000, NumWorkers: 10, DelayMs: 100},
		{NumMessages: 1000, NumWorkers: 20, DelayMs: 100},
		{NumMessages: 1000, NumWorkers: 40, DelayMs: 100},
		{NumMessages: 1000, NumWorkers: 80, DelayMs: 100},
		{NumMessages: 1000, NumWorkers: 160, DelayMs: 100},
		{NumMessages: 1000, NumWorkers: 320, DelayMs: 100},
		{NumMessages: 1000, NumWorkers: 640, DelayMs: 100},
	}

	var results []BenchmarkResult

	for _, test := range tests {
		fmt.Println("========================================")
		fmt.Printf("Benchmarking with %d workers\n", test.NumWorkers)

		log.Println("Publishing messages...")
		err := runProducer(test.NumMessages)
		if err != nil {
			log.Fatalf("Error in producer: %v", err)
		}

		cpuProfileFile, err := createProfileFile(profileDir, "cpu", test.NumWorkers)
		if err != nil {
			log.Fatalf("%v", err)
		}

		memProfileFile, err := createProfileFile(profileDir, "mem", test.NumWorkers)
		if err != nil {
			log.Fatalf("%v", err)
		}

		pprof.StartCPUProfile(cpuProfileFile)
		pprof.WriteHeapProfile(memProfileFile)

		log.Println("Starting consumer...")
		start := time.Now()
		runConsumer(test.NumWorkers, test.NumMessages, test.DelayMs)
		elapsed := time.Since(start)

		results = append(results, BenchmarkResult{
			NumMessages: test.NumMessages,
			NumWorkers:  test.NumWorkers,
			DelayMs:     test.DelayMs,
			Duration:    elapsed,
		})

		pprof.StopCPUProfile()
		cpuProfileFile.Close()

	}

	printBenchmarkResults(results)

}

func printBenchmarkResults(results []BenchmarkResult) {
	fmt.Println("\n========== Benchmark Results ==========")

	fmt.Printf("%-10s %-10s %-10s %-15s\n",
		"Workers", "Messages", "Delay", "Duration")
	fmt.Println(strings.Repeat("-", 100))

	for _, r := range results {
		fmt.Printf("%-10d %-10d %-10d %-15s\n",
			r.NumWorkers,
			r.NumMessages,
			r.DelayMs,
			r.Duration,
		)
	}

	fmt.Println(strings.Repeat("=", 100))
}
