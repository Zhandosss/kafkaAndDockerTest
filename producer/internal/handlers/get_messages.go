package handlers

import (
	"Messaggio/producer/internal/model"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type GetMessageResponse struct {
	Message *model.Message `json:"message"`
}

// getMessage
// @Summary Get message
// @Tags messages
// @Description Get message by id
// @ID get-message
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} GetMessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages/{id} [get]
func (h *Handlers) getMessage(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	log.Info().Msgf("RequestID: %s, Get message request", requestID)

	id := c.Param("id")
	if id == "" {
		log.Error().Msgf("RequestID: %s, ID is required", requestID)
		return c.JSON(http.StatusBadRequest, &ErrorResponse{Message: "id is required"})
	}

	log.Info().Msgf("RequestID: %s, Message ID: %s", requestID, id)

	message, err := h.services.GetMessage(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Error().Err(err).Msgf("RequestID: %s, Message not found", requestID)
			return c.JSON(http.StatusNotFound, &ErrorResponse{Message: "message not found"})
		}
		log.Error().Err(err).Msgf("RequestID: %s, Failed to get message", requestID)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to get message"})
	}

	log.Info().Msgf("RequestID: %s, Message: %+v", requestID, message)
	return c.JSON(http.StatusOK, GetMessageResponse{Message: message})
}

type GetMessagesResponse struct {
	Messages []*model.Message `json:"messages"`
}

// getMessages
// @Summary Get messages
// @Tags messages
// @Description Get all messages
// @ID get-messages
// @Accept json
// @Produce json
// @Success 200 {object} GetMessagesResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages [get]
func (h *Handlers) getMessages(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	log.Info().Msgf("RequestID: %s, Get messages request", requestID)
	messages, err := h.services.GetMessages()
	if err != nil {
		log.Error().Err(err).Msgf("RequestID: %s, Failed to get messages", requestID)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to get messages"})
	}

	log.Info().Msgf("RequestID: %s, returned %d messages", requestID, len(messages))
	return c.JSON(http.StatusOK, GetMessagesResponse{Messages: messages})
}
