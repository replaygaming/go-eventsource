package consumer

import (
	"golang.org/x/net/context"
	"google.golang.org/cloud/pubsub"
	"log"
)

type Consumer struct {
	Subscription *pubsub.Subscription
}

func NewConsumer(topicName string, subscriptionName string) (*Consumer, error) {
	pubsubClient, err := pubsub.NewClient(context.Background(), "emulator-project-id")
	if err != nil {
		log.Fatalf("[FATAL] Could not create PubSub client: %v", err)
		return nil, err
	}

	topic, err := ensureTopic(pubsubClient, topicName)
	if err != nil {
		log.Fatalf("[FATAL] Could not use topic: %v", err)
		return nil, err
	}

	sub, err := ensureSubscription(pubsubClient, topic, subscriptionName)
	if err != nil {
		log.Fatalf("[FATAL] Could not use subscription: %v", err)
		return nil, err
	}

	return &Consumer{Subscription: sub}, nil
}

func ensureTopic(pubsubClient *pubsub.Client, topicName string) (*pubsub.Topic, error) {
	var topic *pubsub.Topic
	topic = pubsubClient.Topic(topicName)
	topicExists, err := topic.Exists(context.Background())
	if err != nil {
		log.Fatalf("[FATAL] Could not verify PubSub topic existence: %v", err)
		return nil, err
	}

	if !topicExists {
		new_topic, err := pubsubClient.NewTopic(context.Background(), topicName)
		if err != nil {
			log.Fatalf("[FATAL] Could not create PubSub topic: %v", err)
			return nil, err
		}
		topic = new_topic
	}

	return topic, nil
}

func ensureSubscription(pubsubClient *pubsub.Client, topic *pubsub.Topic, subscriptionName string) (*pubsub.Subscription, error) {
	var subscription *pubsub.Subscription
	subscription = pubsubClient.Subscription(subscriptionName)
	subscriptionExists, err := subscription.Exists(context.Background())
	if err != nil {
		log.Fatalf("[FATAL] Could not verify PubSub subscription existence: %v", err)
		return nil, err
	}

	if !subscriptionExists {
		new_subscription, err := pubsubClient.NewSubscription(context.Background(), subscriptionName, topic, 0, nil)
		if err != nil {
			log.Fatalf("[FATAL] Could not create PubSub subscription: %v", err)
			return nil, err
		}
		subscription = new_subscription
	}

	return subscription, nil
}

func Consume(consumer *Consumer) (chan *pubsub.Message, error) {
	// Construct the iterator

	channel := make(chan *pubsub.Message)

	go func() {
		it, err := consumer.Subscription.Pull(context.Background())
		if err != nil {
			log.Fatalf("[FATAL] Could not pull message from subscription: %v", err)
			return
		}
		defer it.Stop()

		for {
			msg, err := it.Next()
			if err == pubsub.Done {
				log.Print("[DEBUG] No more messages")
				break
			}
			if err != nil {
				// handle err ...
				log.Fatalf("[FATAL] Error consuming messages: %v", err)
				break
			}

			log.Print("got message: ", string(msg.Data))
			channel <- msg
		}
	}()

	return channel, nil
}
