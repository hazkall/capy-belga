package mq

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewMQ(url string) (*MQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &MQ{
		URL:     url,
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (mq *MQ) DeclareQueues(queueNames []string) error {
	for _, qn := range queueNames {
		_, err := mq.Channel.QueueDeclare(
			qn,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		slog.Info("Queue declared", "queue", qn)
	}
	return nil
}

func (mq *MQ) PublishMessage(message []byte, queueName string) error {
	return mq.Channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

func (mq *MQ) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	return mq.Channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (mq *MQ) Close() {
	if mq.Channel != nil {
		mq.Channel.Close()
	}
	if mq.Conn != nil {
		mq.Conn.Close()
	}
}
