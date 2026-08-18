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

	// create exchange
	err = ch.ExchangeDeclare(
		"lab.exchange", // name
		"direct",       // type
		true,           // durable
		false,          // auto-delete
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	//  create queue
	q, err := ch.QueueDeclare(
		"work.queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Queue -> Exchange binding
	err = ch.QueueBind(
		q.Name,
		"work",
		"lab.exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 50; i++ {
		message := fmt.Sprintf("Message %d", i)

		err = ch.Publish(
			"lab.exchange",
			"work",
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(message),
				DeliveryMode: amqp.Persistent,
			},
		)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Sent:", message)
	}
}
