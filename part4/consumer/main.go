package main

import (
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := ampq.Dial("amqp://guest:qqwwee1@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"ornekuyruk",
		true,
		false,
		false,
		false,
		nil,
	)
	msgs, err := ch.Consume(
		"ornekuyruk", // queue name
		"",           // consumer name
		false,        //auto ack
		false,        // exclusice
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		log.Printf("Received message: %s and lenght:%d", msg.Body, len(msg.Body))

		if err := msg.Ack(false); err != nil {
			log.Println("ACK failed:", err)
		}
	}

}
