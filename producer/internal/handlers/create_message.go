package handlers

import (
	"Messaggio/producer/internal/model"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (h *Handlers) createMessage(c echo.Context) error {
	var message model.Message

	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{Message: "invalid request body"})
	}
	if message.Content == "" {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{Message: "content is required"})
	}

	message.CreateTime = time.Now()

	id, err := h.services.SaveMessage(&message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to save message"})
	}

	message.ID = id

	bytes, err := json.Marshal(message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to marshal message"})
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: "messages",
		Key:   sarama.StringEncoder(id),
		Value: sarama.StringEncoder(bytes),
	}

	_, _, err = h.producer.SendMessage(kafkaMsg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to send message"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}
