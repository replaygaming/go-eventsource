package consumer

import (
	"golang.org/x/net/context"
	"google.golang.org/cloud/pubsub"
	"log"
)

type Consumer struct {
	Subscription *pubsub.Subscription
}

const projectId = "emulator-project-id"

func NewConsumer(topicName string, subscriptionName string) *Consumer {
	pubsubClient, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		log.Fatalf("[FATAL] Could not create PubSub client: %v", err)
	}

	log.Printf("Using Google PubSub with project id: %s", projectId)

	topic := ensureTopic(pubsubClient, topicName)
	sub := ensureSubscription(pubsubClient, topic, subscriptionName)

	return &Consumer{Subscription: sub}
}

func ensureTopic(pubsubClient *pubsub.Client, topicName string) *pubsub.Topic {
	var topic *pubsub.Topic
	topic = pubsubClient.Topic(topicName)
	topicExists, err := topic.Exists(context.Background())
	if err != nil {
		log.Fatalf("[FATAL] Could not verify PubSub topic existence: %v", err)
	}

	if !topicExists {
		log.Println("[INFO] Creating new topic")
		new_topic, err := pubsubClient.NewTopic(context.Background(), topicName)
		if err != nil {
			log.Fatalf("[FATAL] Could not create PubSub topic: %v", err)
		}
		topic = new_topic
	} else {
		log.Println("[INFO] Using existing topic")
	}

	return topic
}

func ensureSubscription(pubsubClient *pubsub.Client, topic *pubsub.Topic, subscriptionName string) *pubsub.Subscription {
	var subscription *pubsub.Subscription
	subscription = pubsubClient.Subscription(subscriptionName)
	subscriptionExists, err := subscription.Exists(context.Background())
	if err != nil {
		log.Fatalf("[FATAL] Could not verify PubSub subscription existence: %v", err)
	}

	if !subscriptionExists {
		log.Println("[INFO] Creating new subscription")
		new_subscription, err := pubsubClient.NewSubscription(context.Background(), subscriptionName, topic, 0, nil)
		if err != nil {
			log.Fatalf("[FATAL] Could not create PubSub subscription: %v", err)
		}
		subscription = new_subscription
	} else {
		log.Println("[INFO] Using existing subscription")
	}

	return subscription
}

func (consumer *Consumer) Consume() (chan *pubsub.Message, error) {
	// Construct the iterator

	channel := make(chan *pubsub.Message)

	go func() {
		it, err := consumer.Subscription.Pull(context.Background())
		if err != nil {
			log.Printf("Could not pull message from subscription: %v", err)
			return
		}
		defer it.Stop()

		for {
			msg, err := it.Next()
			if err == pubsub.Done {
				log.Println("[DEBUG] No more messages")
				break
			}
			if err != nil {
				// handle err ...
				log.Printf("[ERROR] Error consuming messages: %v", err)
				break
			}

			log.Printf("[DEBUG] Got message: ", string(msg.Data))
			channel <- msg
		}
	}()

	return channel, nil
}

func (consumer *Consumer) Remove() {
	log.Printf("[INFO] Removing subscription %s", consumer.Subscription.Name())
	consumer.Subscription.Delete(context.Background())
}
