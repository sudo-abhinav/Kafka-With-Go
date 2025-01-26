package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sudo-abhinav/models"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notification"
)

var ErrUserNotFoundInProducer = errors.New("user not found")

func findUserByID(id int, users []models.User) (models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, ErrUserNotFoundInProducer
}

func getIDFromRequest(r *http.Request, key string) (int, error) {
	idStr := r.FormValue(key)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID from form value %s: %w", key, err)
	}
	return id, nil
}
