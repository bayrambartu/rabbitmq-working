package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:qqwwee1@localhost:5672/")
	failOnError(err, "cannot connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "cannot open a channel")
	defer ch.Close()

	replyQueue, err := ch.QueueDeclare(
		"",
		false, // durable
		false, // auto-delete
		true,
		false,
		nil,
	)
	failOnError(err, "cannot declare a reply queue")

	msgs, err := ch.Consume(
		replyQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "cannot register a consumer")

	corrId := uuid.New().String()

	n := 7
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf(" [x] Sending request: square(%d)", n)

	err = ch.PublishWithContext(ctx,
		"", // default exchange
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       replyQueue.Name,
			Body:          []byte(strconv.Itoa(n)),
		},
	)
	failOnError(err, "cannot publish a message")

	for d := range msgs {
		if d.CorrelationId == corrId {
			fmt.Printf(" [.] received: %s\n", d.Body)
			return
		}
	}
}
