package main

import (
	"context"
	"log"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func square(n int) int {
	return n * n
}
func main() {
	conn, err := amqp.Dial("amqp://guest:qqwwee1@localhost:5672/")
	failOnError(err, "cannot connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "cannot open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "cannot declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "cannot register a consumer")

	log.Printf(" [*] Awaiting RPC requests")

	for d := range msgs {
		n, err := strconv.Atoi(string(d.Body))
		if err != nil {
			log.Printf("Invalid request	: %s", d.Body)
			d.Ack(false)
			continue
		}
		log.Printf(" [.] square(%d)", n)
		result := square(n)

		err = ch.PublishWithContext(
			context.Background(),
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(strconv.Itoa(result)),
			})
		failOnError(err, "cannot publish a message")
		d.Ack(false)
	}

}
