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

	msgs, err := ch.Consume(
		"delay-queue",
		"wkwk",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Fail to register consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			body := "Hello delay queue consumer "
			err := ch.Publish(
				"delay-exchange",
				"",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
					Headers: amqp.Table{
						"x-delay": 10 * 1000,
					},
				},
			)
			failOnError(err, "Fail to publish message")

			log.Printf(" [x] from consumer sent message!")
		}
	}()

	log.Printf(" [*] waiting for messages. to exit pres ctrl c ")
	<-forever
}
