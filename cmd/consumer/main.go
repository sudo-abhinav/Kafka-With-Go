package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sudo-abhinav/cmd/models"

	"golang.org/x/net/context"
	"log"
	"net/http"
)

func main() {
	store := &NotificationStore{
		data: make(UserNotifications),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go SetupConsumerGroup(ctx, store)
	defer cancel()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/notifications/{userID}", func(w http.ResponseWriter, r *http.Request) {
		handleNotifications(w, r, store)
	})

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 started at http://localhost%s\n", ConsumerGroup, ConsumerPort)

	if err := http.ListenAndServe(ConsumerPort, r); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}

func handleNotifications(w http.ResponseWriter, r *http.Request, store *NotificationStore) {
	userID, err := GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	notes := store.Get(userID)
	if len(notes) == 0 {
		response := map[string]interface{}{
			"message":       "No notifications found for user",
			"notifications": []models.Notification{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"notifications": notes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
