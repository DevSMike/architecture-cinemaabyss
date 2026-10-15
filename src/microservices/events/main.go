package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "kafka:9092"
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// consumer
	go startConsumers(ctx, brokers)

	// producer
	writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
	}
	defer writer.Close()

	http.HandleFunc("/api/events/movie", handlePublishMovie)
	http.HandleFunc("/api/events/user", handlePublishUser)
	http.HandleFunc("/api/events/payment", handlePublishPayment)
	http.HandleFunc("/api/events/health", handleHealth)

	log.Println("Events service listening on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
