package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/cloud/pubsub"
	"log"
	"os"
	"strings"
)

type EventSourceMessage struct {
	Event    string   `json:"event"`
	Data     string   `json:"data"`
	Channels []string `json:"channels"`
}

var (
	topicName = flag.String("topic", "eventsource", "Topic name")
	channels  = flag.String("channels", "*", "Channel list (comma separated)")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: client EVENT MESSAGE\n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	pubsubClient, _ := pubsub.NewClient(context.Background(), "emulator-project-id")
	topic := pubsubClient.Topic("eventsource")

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(2)
	}

	event := flag.Arg(0)
	contents := flag.Arg(1)

	channelsList := strings.Split(*channels, ",")

	message := EventSourceMessage{
		Event:    event,
		Data:     contents,
		Channels: channelsList,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Could not parse message: %v", err)
	}

	topic.Publish(context.Background(), &pubsub.Message{
		Data: payload,
	})
}
