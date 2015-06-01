package main

import (
	"log"

	"github.com/streadway/amqp"
)

type consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func (c *consumer) subscribe(amqpURI string) (<-chan amqp.Delivery, error) {
	c.conn = nil
	c.channel = nil
	c.done = make(chan error)

	var err error

	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	go func() {
		log.Fatalf("[FATAL] AMQP closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = c.channel.ExchangeDeclare(
		"es_ex",  // name of the exchange
		"fanout", // type
		true,     // durable
		false,    // delete when complete
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	queue, err := c.channel.QueueDeclare(
		"",    // name of the queue
		true,  // durable
		true,  // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = c.channel.QueueBind(
		queue.Name, // name of the queue
		"es_ex",    // bindingKey
		"es_ex",    // sourceExchange
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}

	return c.channel.Consume(
		queue.Name,    // name
		"eventsource", // consumerTag,
		false,         // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)
}
