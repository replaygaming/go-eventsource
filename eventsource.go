package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/luizbranco/eventsource"
)

var (
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
		ChanSub: eventsource.QueryStringChannels{Name: "channels"},
	}
	stats, err := NewStats(*statsdUrl, *prefix)
	if err == nil {
		server.Stats = stats
	} else {
		log.Printf("[ERROR] %s", err)
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
	log.Printf("[INFO] start port=%s amqp-url=%s statsd-url=%s statsd-prefix=%s"+
		" compression=%t ", *port, *amqpUrl, *statsdUrl, *prefix, *compress)
	err = http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalf("[FATAL] Server %s", err)
	}
}
