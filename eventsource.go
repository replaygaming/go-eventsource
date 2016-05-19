package main

import (
	"log"
	"net/http"
	"os"

	"github.com/replaygaming/eventsource"
	amqp "github.com/streadway/amqp"
)

var (
	environment string
	amqpURL     string
	exchange    string
	prefix      string
	username    string
	password    string
	compress    bool
)

func init() {
	log.Printf("[INFO] Starting")

	environment = os.Getenv("ENV")
	amqpURL = os.Getenv("AMQP_URL")
	exchange = os.Getenv("EXCHANGE")
	prefix = "eventsource"
	username = os.Getenv("EVENTSOURCE_USER")
	password = os.Getenv("EVENTSOURCE_PASSWORD")
	compress = os.Getenv("COMPRESS") != "false"

	log.Printf("[INFO] INIT - environment=%s, AMPQ URL=%s, exchange=%s,"+
		"prefix=%s, compress=%t, username=%s, password=%s", environment, amqpURL,
		exchange, prefix, compress, username, password)
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

func newConsumer(amqpURL string, exchange string) *Consumer {
	// NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error)
	consumer, err := NewConsumer(amqpURL, exchange, "fanout", "", "", "eventsource")
	if err != nil {
		fatal("AMQP consumer failed", err)
	}
	return consumer
}

func consume(consumer *Consumer) <-chan amqp.Delivery {
	messages, err := consumer.Consume()
	if err != nil {
		fatal("AMQP queue failed", err)
	}
	return messages
}

func messageLoop(messages <-chan amqp.Delivery, server *eventsource.Eventsource, consumer *Consumer) {
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

func heartbeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func startServing(server *eventsource.Eventsource) {
	http.HandleFunc("/", heartbeat)
	http.Handle("/subscribe", server)

	log.Printf("[INFO] STARTING - environment=%s, AMPQ URL=%s, exchange=%s,"+
		" stats prefix=%s, compression=%t", environment, amqpURL,
		exchange, prefix, compress)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fatal("Server", err)
	}
}

func main() {
	server := newServerWithMetrics(prefix)
	consumer := newConsumer(amqpURL, exchange)
	messages := consume(consumer)

	go messageLoop(messages, server, consumer)

	startServing(server)
}
