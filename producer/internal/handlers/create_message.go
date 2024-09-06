package handlers

import (
	"Messaggio/producer/internal/model"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type CreateMessageResponse struct {
	ID string `json:"id"`
}

// createMessage
// @Summary Create message
// @Tags messages
// @Description Create message
// @ID create-message
// @Accept json
// @Produce json
// @Param message body model.Message true "Message. Provide only content"
// @Success 201 {object} CreateMessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages [post]
func (h *Handlers) createMessage(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	var message model.Message
	log.Info().Msgf("RequestID: %s, Create message request", requestID)

	if err := c.Bind(&message); err != nil {
		log.Error().Err(err).Msgf("RequestID: %s, Invalid request body", requestID)
		return c.JSON(http.StatusBadRequest, &ErrorResponse{Message: "invalid request body"})
	}
	log.Info().Msgf("RequestID: %s, Message: %+v", requestID, message)

	if message.Content == "" {
		log.Error().Msgf("RequestID: %s, Content is required", requestID)
		return c.JSON(http.StatusBadRequest, &ErrorResponse{Message: "content is required"})
	}

	message.CreateTime = time.Now()

	id, err := h.services.SaveMessage(&message)
	if err != nil {
		log.Error().Err(err).Msgf("RequestID: %s, Failed to save message", requestID)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to save message"})
	}

	message.ID = id

	_, err = json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msgf("RequestID: %s, Failed to marshal message", requestID)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to marshal message"})
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: "messages",
		Key:   sarama.StringEncoder(id),
		Value: sarama.StringEncoder(""),
	}

	_, _, err = h.producer.SendMessage(kafkaMsg)
	if err != nil {
		log.Error().Err(err).Msgf("RequestID: %s, Failed to send message to kafka", requestID)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to send message"})
	}

	log.Info().Msgf("RequestID: %s, Message saved in database and sent to kafka", requestID)
	return c.JSON(http.StatusCreated, CreateMessageResponse{ID: id})
}
