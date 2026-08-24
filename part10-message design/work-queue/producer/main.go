package main

import (
	"fmt"
	"log"
	"time"

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

	queueName := "example-work-queue"
	ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	for i := 0; i < 50; i++ {
		time.Sleep(1 * time.Second)
		body := "Hello World" + fmt.Sprintf("%d", i)

		err = ch.Publish(
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		failOnError(err, "cannot publish a message")

	}

}
