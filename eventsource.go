package main

import (
	"encoding/hex"
	"flag"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/replaygaming/eventsource"
)

var (
	topicName = flag.String("topic",
		fromEnvWithDefault("ES_TOPIC", "eventsource"),
		"Topic name")

	subscriptionName = flag.String("subscription",
		fromEnvWithDefault("ES_SUBSCRIPTION", "eventsource_"+generateSubId()),
		"Subscription name")

	port = flag.String("port",
		fromEnvWithDefault("ES_PORT", "3001"),
		"Eventsource port")

	useMetrics = flag.Bool("metrics", os.Getenv("ES_METRICS") == "true", "Enable metrics")

	metricsPrefix = flag.String("metrics-prefix",
		fromEnvWithDefault("ES_METRICS_PREFIX", "eventsource"),
		"Metrics prefix")

	compress = flag.Bool("compression", os.Getenv("ES_COMPRESSION") == "true", "Enable zlib compression of data")

	verbose = flag.Bool("verbose", os.Getenv("ES_VERBOSE") == "true", "Enable verbose output")
)

func init() {
	flag.Parse()
}

func newServerWithMetrics() *eventsource.Eventsource {
	server := &eventsource.Eventsource{
		ChannelSubscriber: eventsource.QueryStringChannels{Name: "channels"},
	}

	if *useMetrics {
		metrics, err := NewMetrics(*metricsPrefix)
		if err != nil {
			Fatal("Could not instantiate metrics", err)
		}
		server.Metrics = metrics
	}
	server.Start()
	return server
}

func newConsumer() *Consumer {
	return NewConsumer(*topicName, *subscriptionName)
}

func consume(c *Consumer) <-chan Message {
	messages, err := c.Consume()
	if err != nil {
		Fatal("Failed to consume messages", err)
	}
	return messages
}

func messageLoop(messages <-chan Message, server *eventsource.Eventsource, c *Consumer) {
	for m := range messages {
		messageData := m.Data()
		if *verbose {
			Debug("Got message: %s", string(messageData))
		}
		data, channels, err := parseMessage(messageData)
		if err != nil {
			Warn("json conversion failed %s", err)
		} else {
			e := eventsource.DefaultEvent{
				Message:  data,
				Channels: channels,
				Compress: *compress,
			}
			server.Send(e)
		}
		m.Done(true)
	}
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func startServing(server *eventsource.Eventsource) {
	http.HandleFunc("/", heartbeat)
	http.Handle("/subscribe", server)

	Info("EventSource server started")
	Info("Configuration: port=%s subscription=%s topic=%s"+
		" compression=%t metrics=%t", *port, *subscriptionName,
		*topicName, *compress, *useMetrics)
	if *useMetrics {
		Info("Metrics configuration: metrics-prefix=%s", *metricsPrefix)
	}
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		Fatal("Error starting HTTP server: %v", err)
	}
}

func main() {
	server := newServerWithMetrics()
	c := newConsumer()
	messages := consume(c)

	setupSignalHandlers(c)

	go messageLoop(messages, server, c)

	startServing(server)
}

var shuttingDown = false

func setupSignalHandlers(consumer *Consumer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if shuttingDown {
			os.Exit(1)
		}
		Info("Shutting down gracefully. Repeat signal to force shutdown")
		shuttingDown = true
		consumer.Remove()
		os.Exit(0)
	}()
}

func generateSubId() string {
	id := make([]byte, 4)
	todo := len(id)
	offset := 0
	source := rand.NewSource(time.Now().UnixNano())
	for {
		val := int64(source.Int63())
		for i := 0; i < 8; i++ {
			id[offset] = byte(val & 0xff)
			todo--
			if todo == 0 {
				return hex.EncodeToString(id)
			}
			offset++
			val >>= 8
		}
	}
}

func fromEnvWithDefault(varName string, defaultValue string) string {
	value := os.Getenv(varName)
	if value != "" {
		return value
	} else {
		return defaultValue
	}
}
