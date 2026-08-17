package main

import (
	"fmt"
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

	log.Println("Channel created")

	_, err = ch.QueueDeclare("ornekuyruk", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		err = ch.Publish("", "ornekuyruk", false, false, ampq.Publishing{
			// exchange -> "" -> default exchange, routing key -> "ornekuyruk" --> queue name
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("Hello World %d", i)),
		})

	}
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Queue declared")

}
