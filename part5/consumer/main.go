package main

import (
	"fmt"
	"log"
	"time"

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

	q, err := ch.QueueDeclare(
		"work.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Qos is used to this

	msgs, err := ch.Consume(
		q.Name,
		"consumer-1",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("consumer-1 is started")

	for msg := range msgs {
		fmt.Println("Received message:", string(msg.Body))

		time.Sleep(time.Second * 2)

		err := msg.Nack(false, true)
		if err != nil {
			log.Println("Error nacking message:", err)
			continue
		}
		fmt.Println("Message nacked", string(msg.Body))
	}
}
