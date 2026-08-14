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
		"payment",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = ch.QueueDeclare(
		"notification",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	paymentMsgs, err := ch.Consume(
		"payment", // queue name
		"",        // consumer name
		false,     //auto ack
		false,     // exclusice
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		log.Fatal(err)
	}

	notificationMsgs, err := ch.Consume(
		"notification",
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

	// for msg := range msgs {
	// 	log.Printf("Received message: %s", msg.Body)

	// 	log.Println("Sending ack...")
	// 	//time.Sleep(10 * time.Second) // Simulate processing time

	// 	//after processing the message, send an acknowledgment
	// 	err := msg.Ack(false)
	// 	if err != nil {
	// 		log.Println("ACK failed:", err)
	// 	}
	// 	log.Println("ack sent")
	// }

	go func() {
		for msg := range paymentMsgs {
			log.Printf("PAYMENT: %s", msg.Body)

			if err := msg.Ack(false); err != nil {
				log.Println("ACK failed:", err)
			}
		}
	}()

	go func() {
		for msg := range notificationMsgs {
			log.Printf("NOTIFICATION: %s", msg.Body)

			if err := msg.Ack(false); err != nil {
				log.Println("ACK failed:", err)
			}
		}
	}()
	select {}

	// rabbitmq-service start
}
