package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/luizbranco/eventsource"
)

var (
	port     = flag.String("port", "3001", "Eventsource port")
	uri      = flag.String("uri", "amqp://guest:guest@localhost:5672/eventsource", "AMQP URI")
	compress = flag.Bool("compression", false, "Enable zlib compression of data")
)

func init() {
	flag.Parse()
}

func main() {
	server := eventsource.NewServer()
	c := consumer{}
	messages, err := c.subscribe(*uri)
	if err != nil {
		log.Fatalf("[FATAL] AMQP %s", err)
	}

	go func() {
		for m := range messages {
			data, channels, err := parseMessage(m.Body)
			if err != nil {
				log.Printf("[WARN] json conversion failed %s", err)
			} else {
				e := eventsource.Event{
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
	log.Printf("[INFO] start port=%s amqp=%s compression=%t", *port, *uri,
		*compress)
	err = http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalf("[FATAL] Server %s", err)
	}
}
