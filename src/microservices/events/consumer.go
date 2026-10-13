package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

func startConsumers(ctx context.Context, brokers string) {
	var wg sync.WaitGroup

	for _, topic := range Topics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			consume(ctx, brokers, topic)
		}(topic)
	}

	wg.Wait()
	log.Println("all consumers stopped")
}

func consume(ctx context.Context, brokers, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokers},
		Topic:   topic,
		GroupID: "events-service",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Printf("consumer stopped: topic=%s", topic)
				return
			}
			log.Printf("consumer %s: %v", topic, err)
			continue
		}
		handleMessage(topic, msg)
	}
}

func handleMessage(topic string, msg kafka.Message) {
	var ev Event
	if err := json.Unmarshal(msg.Value, &ev); err != nil {
		log.Printf("consumer %s: invalid message: %v", topic, err)
		return
	}
	log.Printf("consumed: topic=%s key=%s value=%s", topic, msg.Key, msg.Value)
}
