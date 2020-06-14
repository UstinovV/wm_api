package main

import (
	"encoding/json"
	"github.com/UstinovV/wm_api/mpsv"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err, "Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err, "Failed to open a channel")
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"mpsv_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err, "Failed to create exchange")
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err, "Failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err, "Failed to register a consumer")
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		"offers",     // routing key
		"mpsv_topic", // exchange
		false,
		nil)
	if err != nil {
		log.Fatal(err, "Failed to bind queue")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			//log.Printf("Received a message: %s", d.Body)
			var offer mpsv.MpsvOffer
			json.Unmarshal(d.Body, &offer)
			log.Println(offer.MpsvId, offer.Title)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
