package main

import (
	"fmt"

	"github.com/google/uuid"
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

	responseQueue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	correlationID := uuid.New().String()

	for i := 0; i < 100; i++ {
		body := fmt.Sprintf("Hello Message %d", i)
		err = ch.Publish(
			"",
			requestQueue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: correlationID,
				ReplyTo:       responseQueue.Name,
				Body:          []byte(body),
			},
		)
		if err != nil {
			panic(err)
		}

	}

	msgs, err := ch.Consume(
		responseQueue.Name,
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

	for msg := range msgs {
		if msg.CorrelationId == correlationID {
			println("Received response:", string(msg.Body))

		}
	}

}
