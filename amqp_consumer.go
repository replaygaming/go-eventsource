package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Consumer wraps the AMQP connection and channel
type Consumer struct {
	conn    *amqp.Connection
	queue   amqp.Queue
	channel *amqp.Channel
	tag     string
	Done    chan error
}

// NewConsumer Creates a new consumer based on the URL and exchange type
func NewConsumer(amqpURI, exchange, exchangeType, queueName, key,
	ctag string) (*Consumer, error) {
	var err error
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		Done:    make(chan error),
	}

	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		log.Fatalf("[FATAL] AMQP closing: %s",
			<-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		true,         // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	c.queue, err = c.channel.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		true,      // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	if err = c.channel.QueueBind(
		c.queue.Name, // name of the queue
		key,          // routing key (ignored on fanout)
		exchange,     // name of the exchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	return c, nil
}

// Consume acts on the consumers queue and returns a Go channel with the queue deliveries
func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	deliveries, err := c.channel.Consume(
		c.queue.Name, // name
		c.tag,        // consumerTag,
		false,        // noAck
		false,        // exclusive
		false,        // noLocal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}
	return deliveries, nil
}
