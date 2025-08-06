package main

import (
	"log"
)

func declareQueue() {
	conn, ch := connect()
	defer conn.Close()
	defer ch.Close()

	queueName := "hello-queue"
	log.Printf("Declaring queue: %s", queueName)
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // auto-delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	log.Printf("Queue declared: %s\n", q.Name)

}
