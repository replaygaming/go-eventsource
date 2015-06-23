package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/luizbranco/eventsource"
)

var (
	env       = flag.String("env", "development", "Environment: development or production")
	port      = flag.String("port", "3001", "Eventsource port")
	amqpUrl   = flag.String("amqp-url", "amqp://guest:guest@localhost:5672/eventsource", "AMQP URL")
	statsdUrl = flag.String("statsd-url", "localhost:8125", "StatsD URL")
	prefix    = flag.String("statsd-prefix", "app.es_go", "StatsD Prefix")
	compress  = flag.Bool("compression", false, "Enable zlib compression of data")
)

func init() {
	flag.Parse()
}

func main() {
	server := &eventsource.Eventsource{
		ChannelSubscriber: eventsource.QueryStringChannels{Name: "channels"},
	}
	metrics, err := NewMetrics(*statsdUrl, *prefix)
	if err == nil && *env == "production" {
		server.Metrics = metrics
	}
	server.Start()

	c := consumer{}
	messages, err := c.subscribe(*amqpUrl)
	if err != nil {
		log.Fatalf("[FATAL] AMQP %s", err)
	}

	go func() {
		for m := range messages {
			data, channels, err := parseMessage(m.Body)
			if err != nil {
				log.Printf("[WARN] json conversion failed %s", err)
			} else {
				e := eventsource.DefaultEvent{
					Message:  data,
					Channels: channels,
					Compress: *compress,
				}
				server.Send(e)
			}
			m.Ack(false)
		}
		c.done <- nil
	}()

	http.Handle("/subscribe", server)
	log.Printf("[INFO] start env=%s port=%s amqp-url=%s statsd-url=%s statsd-prefix=%s"+
		" compression=%t ", *env, *port, *amqpUrl, *statsdUrl, *prefix, *compress)
	err = http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalf("[FATAL] Server %s", err)
	}
}
