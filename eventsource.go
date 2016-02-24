package main

import (
	"log"
	"net/http"
	"os"

	"github.com/luizbranco/eventsource"
	amqpConsumer "github.com/replaygaming/amqp-consumer"
	amqp "github.com/streadway/amqp"
)

var (
	environment string
	port        string
	amqpURL     string
	exchange    string
	prefix      string
	compress    bool
)

func init() {
	environment = os.Getenv("ENV")
	port = os.Getenv("PORT")
	amqpURL = os.Getenv("AMQP_URL")
	exchange = os.Getenv("EXCHANGE")
	prefix = "eventsource"
	compress = os.Getenv("COMPRESS") == "true"
}

func warn(message string, err error) {
	log.Printf("[WARN] %s: %s", message, err)
}

func fatal(message string, err error) {
	log.Fatalf("[FATAL] %s: %s", message, err)
}

func newServerWithMetrics(prefix string) *eventsource.Eventsource {
	server := &eventsource.Eventsource{
		ChannelSubscriber: eventsource.QueryStringChannels{Name: "channels"},
	}
	metrics, err := NewMetrics(prefix)
	if err != nil {
		warn("Metrics proxy creation failed", err)
	} else if environment == "production" {
		server.Metrics = metrics
	}
	server.Start()
	return server
}

func newConsumer(amqpURL string, exchange string) *amqpConsumer.Consumer {
	// NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error)
	consumer, err := amqpConsumer.NewConsumer(amqpURL, exchange, "fanout", "", "", "eventsource")
	if err != nil {
		fatal("AMQP consumer failed", err)
	}
	return consumer
}

func consume(exchange string, consumer *amqpConsumer.Consumer) <-chan amqp.Delivery {
	messages, err := consumer.Consume(exchange)
	if err != nil {
		fatal("AMQP queue failed", err)
	}
	return messages
}

func messageLoop(messages <-chan amqp.Delivery, server *eventsource.Eventsource, consumer *amqpConsumer.Consumer) {
	for message := range messages {
		data, channels, err := parseMessage(message.Body)
		if err != nil {
			warn("JSON conversion failed", err)
		} else {
			event := eventsource.DefaultEvent{
				Message:  data,
				Channels: channels,
				Compress: compress,
			}
			server.Send(event)
		}
		message.Ack(false)
	}
	consumer.Done <- nil
}

func startServing(server *eventsource.Eventsource) {
	http.Handle("/subscribe", server)

	log.Printf("[INFO] STARTING - environment=%s, port=%s, AMPQ URL=%s, exchange=%s,"+
		" stats prefix=%s, compression=%t", environment, port, amqpURL,
		exchange, prefix, compress)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fatal("Server", err)
	}
}

func main() {
	server := newServerWithMetrics(prefix)
	consumer := newConsumer(amqpURL, exchange)
	messages := consume(exchange, consumer)

	go messageLoop(messages, server, consumer)

	startServing(server)
}
