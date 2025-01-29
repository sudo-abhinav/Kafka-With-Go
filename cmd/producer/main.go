package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sudo-abhinav/cmd/models"
	"log"
	"net/http"
)

/*
	func main() {
		http.HandleFunc("/", Sendnotification)
		log.Fatal(http.ListenAndServe(":6900", nil))
	}

	func Sendnotification(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
*/
func main() {
	users := []models.User{
		{Id: 1, Name: "Emma"},
		{Id: 2, Name: "Bruno"},
		{Id: 3, Name: "Rick"},
		{Id: 4, Name: "Lena"},
		{Id: 5, Name: "Michael"},
		{Id: 6, Name: "Jackson"},
		{Id: 7, Name: "chetan"},
		{Id: 8, Name: "Benjamin"},
		{Id: 9, Name: "vikash"},
	}

	producer, err := setupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer func(producer sarama.SyncProducer) {
		err := producer.Close()
		if err != nil {

		}
	}(producer)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/send", sendMessageHandler(producer, users))

	fmt.Printf("Kafka PRODUCER 📨 started at http://localhost%s\n", ProducerPort)

	if err := http.ListenAndServe(ProducerPort, r); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
