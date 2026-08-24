package main

import (
	"context"
	"log"
	"time"
	"fmt"

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

	err = ch.ExchangeDeclare(
		"pub-sub-exchange",   // exchange name
		"fanout", // type
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // no-wait
		nil,
	)
	failOnError(err, "cannot declare exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

for i := 1; i <= 50; i++ {
	body := fmt.Sprintf("System event: user logged in - message %d", i)

	err = ch.PublishWithContext(
		ctx,
		"pub-sub-exchange",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	failOnError(err, "message publish failed")

	log.Printf(" [x] Published: %s", body)
}

}