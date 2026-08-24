package main

import (
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
	ch.Qos(1, 0, false)

	consumer, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for msg := range consumer {
		time.Sleep(1 * time.Second)
		log.Printf("Received a message: %s", msg.Body)

	}
}
