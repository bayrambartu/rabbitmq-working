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

	err = ch.ExchangeDeclare(
		"new-exchange", // name
		"direct",
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
		"payment",
		true,  // durable -> kuyruk kalıcı olsun mu? true -> kalıcı, false -> geçici
		false, // autoDelete -> kuyruk otomatik silinsin mi? true -> silinsin, false -> silinmesin
		false, // exclusive -> eğer bi kuyruk excliszce ise o kutruk o baglantıya özel olusturulur ve sonrasında silinir
		false, // noWait -> bekleme yok mu? true -> bekleme yok, false -> bekleme var
		nil,   //args -> kuyruk için ekstra parametreler
	)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Payment queue declared")

	_, err = ch.QueueDeclare(
		"notification",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Notification queue declared")

	// exchange ile queue'yu baglamak için binding yapıyoruz
	err = ch.QueueBind(
		"payment",         // queue name
		"payment.created", // routing key
		"new-exchange",    // exchange name
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		"notification",         // queue name
		"notification.created", // routing key
		"new-exchange",         // exchange name
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish(
		"new-exchange",    // exchange routing islemini yapar
		"payment.created", // routing key
		false,
		false,
		ampq.Publishing{
			ContentType: "text/plain",
			Body:        []byte("message of world -PAYMENT-"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish(
		"new-exchange",         // exchange routing islemini yapar
		"notification.created", // routing key
		false,
		false,
		ampq.Publishing{
			ContentType: "text/plain",
			Body:        []byte("message of world -NOTIFICATION-"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Messages published")
}
