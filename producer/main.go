package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@xxx:5672/")
	failOnError(err, "Failed to connect to rabbitmq")

	defer conn.Close()

	ch, err := conn.Channel()

	failOnError(err, "Failed to open a channel")

	defer ch.Close()

	body := "Hello delayed message!"

	err = ch.Publish(
		"delay-exchange",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Headers: amqp.Table{
				"x-delay": 5000,
			},
		},
	)

	failOnError(err, "fail to publish message")

	log.Printf("[x] Sent %s ", body)
}
