package main

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	consumerName := os.Getenv("CONSUMER_NAME")

	if consumerName == "" {
		consumerName = "consumer-unknown"
	}

	// RabbitMQ connection
	conn, err := amqp.Dial("amqp://guest:qqwwee1@localhost:5672/")
	if err != nil {
		log.Fatal("RabbitMQ connection error:", err)
	}
	defer conn.Close()

	// Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Channel error:", err)
	}
	defer ch.Close()

	// Queue
	q, err := ch.QueueDeclare(
		"work.queue",
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatal("Queue declare error:", err)
	}

	// --------------------------------------------------
	// QoS
	//
	// Her consumer'ın aynı anda maksimum 2 tane
	// ACK edilmemiş mesajı olabilir.
	// --------------------------------------------------

	err = ch.Qos(
		2,     // prefetch count
		0,     // prefetch size -> unlimited
		false, // global
	)
	if err != nil {
		log.Fatal("QoS error:", err)
	}

	// Consumer
	msgs, err := ch.Consume(
		q.Name,
		consumerName,
		false, // autoAck = false
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		log.Fatal("Consume error:", err)
	}

	fmt.Printf("[%s] Consumer started\n", consumerName)

	// how many messages processed by this consumer
	processed := 0

	for msg := range msgs {

		fmt.Printf(
			"[%s] Received: %s\n",
			consumerName,
			string(msg.Body),
		)

		// --------------------------------------------------
		// İşlem süresini simüle ediyoruz.
		//
		// Consumer A -> 5 saniye
		// Consumer B -> 1 saniye
		// --------------------------------------------------

		sleepTime := 5 * time.Second

		if consumerName == "consumer-B" {
			sleepTime = 1 * time.Second
		}

		fmt.Printf(
			"[%s] Processing for %v...\n",
			consumerName,
			sleepTime,
		)

		time.Sleep(sleepTime)

		// --------------------------------------------------
		// İşlem başarılı.
		// ACK gönderiyoruz.
		// --------------------------------------------------

		err := msg.Ack(false)
		if err != nil {
			log.Printf(
				"[%s] ACK error: %v\n",
				consumerName,
				err,
			)
			continue
		}

		processed++

		fmt.Printf(
			"[%s] ACK: %s | processed=%d\n",
			consumerName,
			string(msg.Body),
			processed,
		)
	}
}
