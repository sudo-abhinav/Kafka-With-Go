package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	ConsumerGroup      = "notifications-group"
	ConsumerTopic      = "notification"
	ConsumerPort       = ":8081"
	KafkaServerAddress = "localhost:9092"
)

var ErrNoMessagesFound = errors.New("no messages found")

func GetUserIDFromRequest(r *http.Request) (string, error) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		return "", ErrNoMessagesFound
	}
	return userID, nil
}
