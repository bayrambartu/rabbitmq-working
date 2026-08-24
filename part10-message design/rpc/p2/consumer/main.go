package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:qqwwee1@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	requestQueueName := "example-reques-response-queue"
	requestQueue, err := ch.QueueDeclare(
		requestQueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}
	msg, err := ch.Consume(
		requestQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for msgs := range msg {

		fmt.Printf("Received a message: %s\n", msgs.Body)

		responseMessage := fmt.Sprintf(
			"Received a message: %s",
			msgs.Body,
		)

		err := ch.Publish(
			"",
			msgs.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: msgs.CorrelationId,
				Body:          []byte(responseMessage),
			},
		)

		if err != nil {
			panic(err)
		}
	}

}
