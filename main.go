package main

import (
	"fmt"
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {

	// canhge the password -> rabbitmqctl change_password guest qqwwee1
	conn, err := ampq.Dial("amqp://guest:qqwwee1@localhost:5672/")
	// rabbitmq'ya bir conneciton açıyor --> conn,err := ampq.Dial("...")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	log.Println("Channel created")

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
	log.Println("Queue declared")

	ch.Publish(
		"",
		"alerts",
		false,
		false,
		ampq.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello from producer!, Merhaba RabbitMQ"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Message published")
}
