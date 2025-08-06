package main

import (
	"fmt"
	"log"
	"math/rand"

	amqp "github.com/rabbitmq/amqp091-go"
)

func runProducer(numMessages int) error {
	conn, ch := connect()
	defer conn.Close()
	defer ch.Close()

	for i := range numMessages {
		body := generateRandomMessage(i)
		err := ch.Publish(
			"",            // default exchange
			"hello-queue", // routing key = queue name
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		if err != nil {
			return err
		}
		log.Printf("Published message %d: %s", i, body)
	}

	log.Printf("Finished publishing %d messages", numMessages)
	return nil
}

func generateRandomMessage(i int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 10)
	for j := range b {
		b[j] = chars[rand.Intn(len(chars))]
	}
	return fmt.Sprintf("Msg #%d: %s", i, string(b))
}
