package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func connect() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatalf("Failed to open channel: %v", err)
	}

	log.Println("Connected to RabbitMQ")
	return conn, ch
}
