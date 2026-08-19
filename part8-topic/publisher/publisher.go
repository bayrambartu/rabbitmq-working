package main

import (
	"fmt"
	"time"

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

	for i := 0; i < 100; i++ {
		time.Sleep(2 * time.Millisecond)

		message := []byte("Hello World!" + fmt.Sprintf("%d", i))

		var topic string
		fmt.Print("Enter the topic format for sending the message: ")
		fmt.Scanln(&topic)

		ch.Publish(
			"topic-exchange-example", // exchange
			topic,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        message,
			},
		)
	}

}
