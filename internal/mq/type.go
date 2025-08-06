package mq

import amqp "github.com/rabbitmq/amqp091-go"

type MQ struct {
	URL     string
	Conn    *amqp.Connection
	Channel *amqp.Channel
}
