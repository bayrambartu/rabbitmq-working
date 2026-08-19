package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, _ := amqp.Dial("amqp://guest:qqwwee1@localhost:5672/")

	ch, _ := conn.Channel()

	ch.ExchangeDeclare(
		"fanout_exchange_example",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	var qname string
	fmt.Print("Enter queue name: ")
	fmt.Scanln(&qname)

	ch.QueueDeclare(qname, true, false, false, false, nil)

	ch.QueueBind(
		qname,
		"",
		"fanout_exchange_example",
		false,
		nil,
	)
	fmt.Println("Queue created:", qname)

	messages, _ := ch.Consume(
		qname,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for msg := range messages {
		fmt.Println("Received message:", string(msg.Body))
	}

}
