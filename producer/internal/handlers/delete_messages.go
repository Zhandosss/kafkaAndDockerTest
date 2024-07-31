package handlers

import (
	"Messaggio/producer/internal/model"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

// deleteMessages
// @Summary Delete messages
// @Tags messages
// @Description Delete all messages
// @ID delete-messages
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Router /messages [delete]
func (h *Handlers) deleteMessages(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	log.Info().Msgf("RequestID: %s, Delete messages request", requestID)
	err := h.services.DeleteMessages()
	if err != nil {
		log.Error().Err(err).Msgf("RequestID: %s, Error deleting messages", requestID)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error deleting messages"})
	}

	log.Info().Msgf("RequestID: %s, Messages deleted", requestID)
	return c.NoContent(http.StatusNoContent)
}

// deleteMessage
// @Summary Delete message
// @Tags messages
// @Description Delete message by id
// @ID delete-message
// @Param id path string true "Message ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages/{id} [delete]
func (h *Handlers) deleteMessage(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	log.Info().Msgf("RequestID: %s, Delete message request", requestID)

	id := c.Param("id")
	if id == "" {
		log.Error().Msgf("RequestID: %s, ID is required", requestID)
		return c.JSON(http.StatusBadRequest, &ErrorResponse{Message: "id is required"})
	}
	log.Info().Msgf("RequestID: %s, Message ID for deletion: %s", requestID, id)

	err := h.services.DeleteMessage(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Error().Err(err).Msgf("RequestID: %s, Message not found", requestID)
			return c.JSON(http.StatusNotFound, &ErrorResponse{Message: "message not found"})
		}
		log.Error().Err(err).Msgf("RequestID: %s, Error deleting message", requestID)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "error deleting message"})
	}

	log.Info().Msgf("RequestID: %s, Message deleted", requestID)
	return c.NoContent(http.StatusNoContent)
}
