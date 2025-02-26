package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/sudo-abhinav/cmd/models"
	"net/http"
	"strconv"
)

func SendKafkaMessage(producer sarama.SyncProducer, users []models.User, r *http.Request, fromID, toID int) error {
	message := r.FormValue("message")

	fromUser, err := FindUserByID(fromID, users)
	if err != nil {
		return err
	}

	toUser, err := FindUserByID(toID, users)
	if err != nil {
		return err
	}

	notification := models.Notification{
		From:    fromUser,
		To:      toUser,
		Message: message,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(strconv.Itoa(toUser.Id)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func sendMessageHandler(producer sarama.SyncProducer, users []models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r)
		fromID, err := GetIDFromRequest(r, "FromId")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		toID, err := GetIDFromRequest(r, "toId")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = SendKafkaMessage(producer, users, r, fromID, toID)
		if errors.Is(err, ErrUserNotFoundInProducer) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Notification sent successfully!"))
		if err != nil {
			return
		}
	}
}

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	//Initializes a new default configuration for Kafka. Think of it as setting up the parameters before connecting to the broker.
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}
