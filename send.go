package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/UstinovV/wm_api/mpsv"
	"io"
	"log"
	"github.com/streadway/amqp"
	"os"
)
func main() {
	//docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Error connecting to RMQ :", err)
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

	//q, err := ch.QueueDeclare(
	//	"hello", // name
	//	false,   // durable
	//	false,   // delete when unused
	//	false,   // exclusive
	//	false,   // no-wait
	//	nil,     // arguments
	//)
	//if err != nil {
	//	log.Fatal(err, "Failed to declare a queue")
	//}
	//todo: receive file from queue
	xmlFile, err := os.Open("example.xml")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err, "Failed to open file")
	}

	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	offer := mpsv.MpsvOffer{}
	for {

		token, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			break

		} else if tokenErr == io.EOF {
			break
		}
		if token == nil {
			fmt.Println("t is nil break")
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "VOLNEMISTO":
				for _, attr := range tok.Attr {
					if attr.Name.Local == "uid" {
						offer.MpsvId = attr.Value
					}
					if attr.Name.Local == "zmena" {
						//offer.CreatedAt = attr.Value
					}
				}
			case "PROFESE":
				for _, attr := range tok.Attr {
					if attr.Name.Local == "nazev" {
						//fmt.Println(attr.Value)
						offer.Title = attr.Value
					}
				}
			case "POZNAMKA":
				decoder.DecodeElement(&offer.Content, &tok)

			}

		case xml.EndElement:
			if tok.Name.Local == "VOLNEMISTO" {
				messageBody, err := json.Marshal(offer)
				if err != nil {
					fmt.Println("Cant send message")
				}
				err = ch.Publish(
					"mpsv_topic",     // exchange
					"offers", // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing {
						ContentType: "application/json",
						Body:        messageBody,
					})
				if err != nil {
					log.Fatal(err, "Failed to publish a message")
				}
				fmt.Println("Sended")
			}
		}
	}

}