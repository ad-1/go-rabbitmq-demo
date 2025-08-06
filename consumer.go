package main

import (
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func runConsumer(numWorkers, numMessages, delayMs int) {

	log.Printf("Running consumer with %d workers and %d ms delay", numWorkers, delayMs)

	conn, ch := connect()
	defer conn.Close()
	defer ch.Close()

	// QoS: only send this many unacked messages at once
	ch.Qos(numWorkers, 0, false)

	msgs, err := ch.Consume(
		"hello-queue",
		"",
		true,  // auto-ack
		false, // not exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	// Internal Go channel between your message consumer and your worker goroutines
	jobChan := make(chan amqp.Delivery)
	var wg sync.WaitGroup

	// Launch workers
	for i := range numWorkers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, delayMs, jobChan) // workers will exit when they finish their current job
		}(i)
	}

	log.Println("Waiting for messages...")
	processedMessages := 0
	for msg := range msgs {
		jobChan <- msg
		processedMessages++
		if processedMessages >= numMessages {
			log.Printf("Processed %d messages, stopping consumer", processedMessages)
			close(jobChan) // Close the job channel to signal workers to stop
			break
		}
	}

	wg.Wait() // Wait for all workers to finish

}

func worker(id, delayMs int, jobs <-chan amqp.Delivery) {
	for job := range jobs {
		log.Printf("Worker %d processing message: %s", id, job.Body)
		time.Sleep(time.Duration(delayMs) * time.Millisecond) // Simulate work
		log.Printf("Worker %d done processing message: %s", id, job.Body)
	}
}
