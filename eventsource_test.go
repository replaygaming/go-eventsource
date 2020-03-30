package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)

var _ = func() bool {
	os.Setenv("PUBSUB_EMULATOR_HOST", fromEnvWithDefault("PUBSUB_EMULATOR_HOST", "pubsub-emulator:8538"))
	testing.Init()
	return true
}()

type EventSourceMessage struct {
	Event    string   `json:"event"`
	Data     string   `json:"data"`
	Channels []string `json:"channels"`
}

var channels = "*"

func sendMessageFromNewClient(event string, contents string) {
	pubsubClient, _ := pubsub.NewClient(context.Background(), "emulator-project-id")
	topic := pubsubClient.Topic("eventsource")

	channelsList := strings.Split(channels, ",")

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

func makeRequest(url string, ch chan<- string) {
	data := make([]byte, 10*1024)

	res, err := http.Get(url)
	if err != nil {
		close(ch)
		fmt.Errorf("Cannot communicate with pubsub")
		os.Exit(1)
	}
	defer res.Body.Close()
	defer close(ch) //don't forget to close the channel as well

	for n, err := res.Body.Read(data); err == nil; n, err = res.Body.Read(data) {
		ch <- string(data[:n])
	}
}

func TestMain(m *testing.T) {
	go main()

	time.Sleep(1 * time.Second)
	ch := make(chan string)
	go makeRequest("http://localhost/subscribe?channels=*,100", ch)
	go makeRequest("http://localhost/subscribe?channels=*,200", ch)
	go makeRequest("http://localhost/subscribe?channels=*,300", ch)

	go func(m *testing.T) {
		time.Sleep(10 * time.Second)
		m.Errorf("Timeout: Did not get all messages!")
	}(m)

	for i := 0; i < 10; i++ {
		go sendMessageFromNewClient("test", "msg"+strconv.Itoa(i))
	}

	done := 0
	for v := range ch {
		for _, sub := range strings.Split(v, "\n") {
			if strings.Contains(sub, "data") {
				done++
			}
		}
		if done == 30 {
			break
		}
	}
}
