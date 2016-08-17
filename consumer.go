package main

import (
	"golang.org/x/net/context"
	"google.golang.org/cloud/pubsub"
	"os"
)

type Consumer struct {
	Subscription *pubsub.Subscription
}

type Message interface {
	Data() []byte
	Done(ack bool)
}

type googlePubSubMessage struct {
	OriginalMessage *pubsub.Message
}

func (m *googlePubSubMessage) Data() []byte {
	return m.OriginalMessage.Data
}

func (m *googlePubSubMessage) Done(ack bool) {
	m.OriginalMessage.Done(ack)
}

var defaultProjectId = "emulator-project-id"

func NewConsumer(topicName string, subscriptionName string) *Consumer {
	projectId := os.Getenv("ES_GOOGLE_PROJECT_ID")
	if projectId == "" {
		projectId = defaultProjectId
	}

	pubsubClient, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		Fatal("Could not create PubSub client: %v", err)
	}

	Info("Using Google PubSub with project id: %s", projectId)

	topic := ensureTopic(pubsubClient, topicName)
	sub := ensureSubscription(pubsubClient, topic, subscriptionName)

	return &Consumer{Subscription: sub}
}

func ensureTopic(pubsubClient *pubsub.Client, topicName string) *pubsub.Topic {
	var topic *pubsub.Topic
	topic = pubsubClient.Topic(topicName)
	topicExists, err := topic.Exists(context.Background())
	if err != nil {
		Fatal("Could not verify PubSub topic existence: %v", err)
	}

	if !topicExists {
		Info("Creating new topic")
		new_topic, err := pubsubClient.NewTopic(context.Background(), topicName)
		if err != nil {
			Fatal("Could not create PubSub topic: %v", err)
		}
		topic = new_topic
	} else {
		Info("Using existing topic")
	}

	return topic
}

func ensureSubscription(pubsubClient *pubsub.Client, topic *pubsub.Topic, subscriptionName string) *pubsub.Subscription {
	var subscription *pubsub.Subscription
	subscription = pubsubClient.Subscription(subscriptionName)
	subscriptionExists, err := subscription.Exists(context.Background())
	if err != nil {
		Fatal("Could not verify PubSub subscription existence: %v", err)
	}

	if !subscriptionExists {
		Info("Creating new subscription")
		new_subscription, err := pubsubClient.NewSubscription(context.Background(), subscriptionName, topic, 0, nil)
		if err != nil {
			Fatal("Could not create PubSub subscription: %v", err)
		}
		subscription = new_subscription
	} else {
		Info("Using existing subscription")
	}

	return subscription
}

func (consumer *Consumer) Consume() (chan Message, error) {
	channel := make(chan Message)

	go func() {
		it, err := consumer.Subscription.Pull(context.Background())
		if err != nil {
			Warn("Could not pull message from subscription: %v", err)
			return
		}
		defer it.Stop()

		for {
			msg, err := it.Next()
			if err == pubsub.Done {
				break
			}
			if err != nil {
				Warn("Error consuming messages: %v", err)
				break
			}

			wrappedMsg := &googlePubSubMessage{OriginalMessage: msg}

			channel <- wrappedMsg
		}
	}()

	return channel, nil
}

func (consumer *Consumer) Remove() {
	Info("Removing subscription %s", consumer.Subscription.Name())
	consumer.Subscription.Delete(context.Background())
}
