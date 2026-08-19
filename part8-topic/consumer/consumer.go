package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, _ := amqp.Dial("amqp://guest:qqwwee1@localhost:5672/")

	ch, _ := conn.Channel()

	ch.ExchangeDeclare(
		"topic-exchange-example",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	var topic string
	fmt.Print("Listening topic format: ")
	fmt.Scanln(&topic)

	queue, _ := ch.QueueDeclare(
		"topic-queue-example",
		true,
		false,
		false,
		false,
		nil,
	)

	ch.QueueBind(
		queue.Name,
		topic, // key string
		"topic-exchange-example",
		false,
		nil,
	)

	msgs, _ := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs {
		fmt.Printf(
			"Message received: %s | Routing Key: %s\n",
			string(msg.Body),
			msg.RoutingKey,
		)
	}

}
