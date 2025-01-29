package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"time"
)

// Kafka configuration
const (
	kafkaTopic  = "notification"
	kafkaBroker = "localhost:9092"
	groupID     = "notification-group"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Route to produce notifications
	r.Post("/send", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("message")
		if msg == "" {
			http.Error(w, "Message is required", http.StatusBadRequest)
			return
		}
		if err := produceMessage(msg); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sent successfully!"))
	})

	// Route to consume notifications
	r.Get("/listen", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientChan := make(chan string, 100)

		// Start consuming Kafka messages
		go func() {
			if err := consumeMessages(ctx, clientChan); err != nil {
				log.Printf("Error in Kafka consumer: %v", err)
			}
			close(clientChan)
		}()

		// Stream messages to client using SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		for msg := range clientChan {
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush()
		}
	})

	log.Println("Server running on :6900")
	log.Fatal(http.ListenAndServe(":6900", r))
}

// produceMessage sends a message to the Kafka topic
func produceMessage(message string) error {
	writer := kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	return writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
}

// consumeMessages reads messages from the Kafka topic
func consumeMessages(ctx context.Context, clientChan chan<- string) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   "notification",
		GroupID: groupID,
		MaxWait: 100 * time.Millisecond,
	})
	defer r.Close()

	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			return err
		}
		clientChan <- string(msg.Value)
	}
}
