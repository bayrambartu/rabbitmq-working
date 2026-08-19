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
		"fanout", // type
		true,
		false,
		false,
		false,
		nil,
	)

	for i := 0; i < 10; i++ {
		messages := []byte("Hello World " + fmt.Sprint(i))
		ch.Publish(
			"fanout_exchange_example", // exchange
			"",                        // The routing key is not used, so it should be empty.
			false,                     // mandatory
			false,                     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        messages,
			})
	}

	fmt.Println("all messages sent")

}
