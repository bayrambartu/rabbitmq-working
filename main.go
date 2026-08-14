package main

import (
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	// declare conneciton
	conn, err := ampq.Dial("amqp://guest:qqwwee1@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// declare channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// declare queue
	_, err = ch.QueueDeclare(
		"example-queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	message := "Hello, RabbitMQ"
	byteMessage := []byte(message)

	err = ch.Publish("", "example-queue", false, false, ampq.Publishing{
		ContentType: "text/plain",
		Body:        byteMessage,
	})
	if err != nil {
		log.Fatalf("cannot publish message: %w", err)
	}

	select {}

}
