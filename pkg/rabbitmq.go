package pkg

import (
	"github.com/streadway/amqp"
)

var rabbitMQ *amqp.Connection

func InitRabbitMQ() {
	var err error
	rabbitMQ, err = amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		panic("Failed to connect to RabbitMQ")
	}
}

func CloseRabbitMQ() {
	if rabbitMQ != nil {
		rabbitMQ.Close()
	}
}
