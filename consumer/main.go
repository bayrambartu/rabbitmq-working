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
		"alerts",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		"alerts",
		"",
		false, //manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)

		log.Println("Sending ack...")
		//time.Sleep(10 * time.Second) // Simulate processing time

		// after processing the message, send an acknowledgment
		err := msg.Ack(false) // multiple=false , requeue=true
		if err != nil {
			log.Println("ACK failed:", err)
		}
		log.Println("ack sent")
	}
}

// rabbitmq-service start
