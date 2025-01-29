package main

import (
	"errors"
	"fmt"
	"github.com/sudo-abhinav/cmd/models"
	"net/http"
	"strconv"
)

const (
	ProducerPort       = ":8088"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notification"
)

var ErrUserNotFoundInProducer = errors.New("user not found")

func FindUserByID(id int, users []models.User) (models.User, error) {
	for _, user := range users {
		if user.Id == id {
			return user, nil
		}
	}
	return models.User{}, ErrUserNotFoundInProducer
}

func GetIDFromRequest(r *http.Request, key string) (int, error) {
	idStr := r.FormValue(key)
	id, err := strconv.Atoi(idStr)

	if err != nil {

		return 0, fmt.Errorf("failed to parse ID from form value %s: %w", key, err)
	}
	return id, nil
}
