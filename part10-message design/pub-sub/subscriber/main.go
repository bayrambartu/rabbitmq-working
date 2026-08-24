package main

import (
	"log"
	"time"
	

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:qqwwee1@localhost:5672/")
	failOnError(err, "cannot connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "cannot open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"pub-sub-exchange",
		"fanout",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "cannot declare exchange")

	q, err := ch.QueueDeclare(
		"", // queue name
		false, // non-durable
		false, // auto-delete
		true, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "cannot declare queue")

	err = ch.QueueBind(
		q.Name,
		"",
		"pub-sub-exchange",
		false,
		nil,
	)
	failOnError(err, "cannot bind queue to exchange")

	// Aynı anda sadece 1 unacked mesaj
	err = ch.Qos(1, 0, false)
	failOnError(err, "cannot set QoS")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // autoAck
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	failOnError(err, "cannot register consumer")

	log.Printf(" [*] Waiting for messages...")

	for d := range msgs {
		log.Printf(" [x] Received: %s", d.Body)

		time.Sleep(2 * time.Second)

		err := d.Ack(false)
		if err != nil {
			log.Println("ACK failed:", err)
			continue
		}

		log.Printf(" [✓] ACK: %s", d.Body)
	}
}
