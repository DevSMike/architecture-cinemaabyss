package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func publishEvent(w http.ResponseWriter, r *http.Request, topic string, ev Event) {
	value, err := json.Marshal(ev)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to marshal event")
		return
	}

	err = writer.WriteMessages(r.Context(), kafka.Message{
		Topic: topic,
		Key:   []byte(ev.ID),
		Value: value,
	})
	if err != nil {
		log.Printf("failed to write message to %s: %v", topic, err)
		writeError(w, http.StatusInternalServerError, "failed to publish event")
		return
	}

	log.Printf("sent event: topic=%s id=%s", topic, ev.ID)
	writeJSON(w, http.StatusCreated, EventResponse{Status: "success", Event: ev})
}

func handlePublishMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var event MovieEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	publishEvent(w, r, TopicMovies, Event{
		ID:        fmt.Sprintf("movie-%d-%s", event.MovieID, event.Action),
		Type:      "movie",
		Timestamp: time.Now().UTC(),
		Payload:   event,
	})
}

func handlePublishUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var event UserEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	publishEvent(w, r, TopicUsers, Event{
		ID:        fmt.Sprintf("user-%d-%s", event.UserID, event.Action),
		Type:      "user",
		Timestamp: event.Timestamp,
		Payload:   event,
	})
}

func handlePublishPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var event PaymentEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	publishEvent(w, r, TopicPayments, Event{
		ID:        fmt.Sprintf("payment-%d-%s", event.PaymentID, event.Status),
		Type:      "payment",
		Timestamp: event.Timestamp,
		Payload:   event,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"status": true})
}
