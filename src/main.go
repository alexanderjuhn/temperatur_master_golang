package main

import (
	"log"

	dc "backend/databaseConnector"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	dc.ReadConfig()
	receive()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func receive() {
	conn, err := amqp.Dial("amqp://" + dc.RabbitConnection.UserRabbit + ":" + dc.RabbitConnection.PasswordRabbit + "@" + dc.RabbitConnection.HostRabbit + ":" + dc.RabbitConnection.PortRabbit + "/")

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		dc.RabbitConnection.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			dc.ProcessValue(string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
