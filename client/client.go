package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"
	"io/ioutil"
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
	pubsubClient, _ := newPubSubClient()
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

var defaultProjectId = "emulator-project-id"

func newPubSubClient() (*pubsub.Client, error) {
	ctx := context.Background()
	projectId := os.Getenv("ES_PUBSUB_PROJECT_ID")
	if projectId == "" {
		projectId = defaultProjectId
	}
	log.Printf("Using Google PubSub with project id: %s", projectId)
	var client *pubsub.Client
	var err error

	// Create a new client with token
	keyfilePath := os.Getenv("ES_PUBSUB_KEYFILE")
	if keyfilePath != "" {
		log.Printf("Using keyfile to authenticate")
		jsonKey, err := ioutil.ReadFile(keyfilePath)
		if err != nil {
			return nil, err
		}
		conf, err := google.JWTConfigFromJSON(
			jsonKey,
			pubsub.ScopeCloudPlatform,
			pubsub.ScopePubSub,
		)

		if err != nil {
			return nil, err
		}
		tokenSource := conf.TokenSource(ctx)
		client, err = pubsub.NewClient(ctx, projectId, cloud.WithTokenSource(tokenSource))
	} else {
		// Create client without token
		client, err = pubsub.NewClient(ctx, projectId)
	}

	return client, err
}
