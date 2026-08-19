package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:qqwwee1@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	if err = ch.ExchangeDeclare("direct-exchange", "direct", true, false, false, false, nil); err != nil {
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"direct-queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// conneciton between queue and exchange
	err = ch.QueueBind(
		q.Name,
		"direct-key",
		"direct-exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Consumer is running...")
	fmt.Println("routing key: order.created")

	// listen for messages
	for msg := range messages {
		fmt.Printf(
			"Mesaj geldi: %s | Routing Key: %s\n",
			string(msg.Body),
			msg.RoutingKey,
		)
	}
}
